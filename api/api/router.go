package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"

	"github.com/gojek/mlp/api/middleware"
	"github.com/gojek/mlp/api/models"
	"github.com/gojek/mlp/api/pkg/authz/enforcer"
	"github.com/gojek/mlp/api/pkg/instrumentation/newrelic"
	"github.com/gojek/mlp/api/pkg/pipelines"
	"github.com/gojek/mlp/api/service"
	"github.com/gojek/mlp/api/validation"
)

type AppContext struct {
	AccountService     service.AccountService
	ApplicationService service.ApplicationService
	ProjectsService    service.ProjectsService
	SecretService      service.SecretService
	UsersService       service.UsersService
	PipelineService    pipelines.PipelineInterface

	AuthorizationEnabled bool
	Enforcer             enforcer.Enforcer

	GitlabEnabled bool
	GitlabService service.GitlabService
}

// type ApiHandler func(r *http.Request, vars map[string]string, body interface{}) *ApiResponse
type ApiHandler func(r *http.Request, vars map[string]string, body interface{}) *ApiResponse

type Route struct {
	method  string
	path    string
	body    interface{}
	handler ApiHandler
	name    string
}

func (route Route) HandlerFunc(validate *validator.Validate) http.HandlerFunc {
	var bodyType reflect.Type
	if route.body != nil {
		bodyType = reflect.TypeOf(route.body)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		for k, v := range r.URL.Query() {
			if len(v) > 0 {
				vars[k] = v[0]
			}
		}

		response := func() *ApiResponse {
			vars["user"] = r.Header.Get("User-Email")
			var body interface{} = nil

			if bodyType != nil {
				body = reflect.New(bodyType).Interface()
				if err := json.NewDecoder(r.Body).Decode(body); err != nil {
					return BadRequest(fmt.Sprintf("Failed to deserialize request body: %s", err.Error()))
				} else if err := validate.Struct(body); err != nil {
					errMessage := err.(validator.ValidationErrors)[0].Translate(validation.EN)
					return BadRequest(errMessage)
				}
			}
			return route.handler(r, vars, body)
		}()

		response.WriteTo(w)
	}
}

func NewRouter(appCtx AppContext) *mux.Router {
	validator := validation.NewValidator()
	applicationsController := ApplicationsController{&appCtx}
	usersController := UsersController{&appCtx}
	projectsController := ProjectsController{&appCtx}
	secretController := SecretsController{&appCtx}
	pipelinesController := PipelinesController{&appCtx}

	routes := []Route{
		// Applications API
		{http.MethodGet, "/applications", nil, applicationsController.ListApplications, "ListApplications"},

		// Users API
		{http.MethodGet, "/users/token/generate", nil, usersController.GenerateToken, "GenerateToken"},
		{http.MethodGet, "/users/authorize", nil, usersController.AuthorizeUser, "AuthorizeUser"},
		{http.MethodGet, "/users/token/retrieve", nil, usersController.RetrieveToken, "RetrieveToken"},

		// Projects API
		{http.MethodGet, "/projects/{project_id:[0-9]+}", nil, projectsController.GetProject, "GetProject"},
		{http.MethodGet, "/projects", nil, projectsController.ListProjects, "ListProjects"},
		{http.MethodPost, "/projects", models.Project{}, projectsController.CreateProject, "CreateProject"},
		{http.MethodPut, "/projects/{project_id:[0-9]+}", models.Project{}, projectsController.UpdateProject, "UpdateProject"},

		// Secret Management API
		{http.MethodGet, "/projects/{project_id:[0-9]+}/secrets", nil, secretController.ListSecret, "ListSecret"},
		{http.MethodPost, "/projects/{project_id:[0-9]+}/secrets", models.Secret{}, secretController.CreateSecret, "CreateSecret"},
		{http.MethodPatch, "/projects/{project_id:[0-9]+}/secrets/{secret_id}", models.Secret{}, secretController.UpdateSecret, "UpdateSecret"},
		{http.MethodDelete, "/projects/{project_id:[0-9]+}/secrets/{secret_id}", nil, secretController.DeleteSecret, "DeleteSecret"},

		// Pipelines API
		{http.MethodGet, "/pipelines", nil, pipelinesController.ListPipelines, "ListPipelines"},
	}

	if appCtx.GitlabEnabled {
		userRoutes := []Route{
			{http.MethodGet, "/users/token/generate", nil, usersController.GenerateToken, "GenerateToken"},
			{http.MethodGet, "/users/authorize", nil, usersController.AuthorizeUser, "AuthorizeUser"},
			{http.MethodGet, "/users/token/retrieve", nil, usersController.RetrieveToken, "RetrieveToken"},
		}
		routes = append(routes, userRoutes...)
	}

	var authzMiddleware *middleware.Authorizer
	var projCreationMiddleware *middleware.ProjectCreation

	if appCtx.AuthorizationEnabled {
		authzMiddleware = middleware.NewAuthorizer(appCtx.Enforcer)
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		_, handler := newrelic.WrapHandle(r.name, r.HandlerFunc(validator))

		if r.name == "CreateProject" {
			handler = projCreationMiddleware.ProjectCreationMiddleware(handler)
		}

		if appCtx.AuthorizationEnabled {
			handler = authzMiddleware.AuthorizationMiddleware(handler)
		}

		router.Name(r.name).
			Methods(r.method).
			Path(r.path).
			Handler(handler)
	}

	return router
}
