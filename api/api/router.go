package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/middleware"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
	"github.com/caraml-dev/mlp/api/pkg/instrumentation/newrelic"
	"github.com/caraml-dev/mlp/api/repository"
	"github.com/caraml-dev/mlp/api/service"
	"github.com/caraml-dev/mlp/api/validation"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Controller interface {
	Routes() []Route
}

type AppContext struct {
	ApplicationService service.ApplicationService
	ProjectsService    service.ProjectsService
	SecretService      service.SecretService

	AuthorizationEnabled bool
	Enforcer             enforcer.Enforcer
}

func NewAppContext(db *gorm.DB, cfg *config.Config) (ctx *AppContext, err error) {
	var authEnforcer enforcer.Enforcer
	if cfg.Authorization.Enabled {
		authEnforcer, err = enforcer.NewEnforcerBuilder().
			URL(cfg.Authorization.KetoServerURL).
			Product("mlp").
			Build()

		if err != nil {
			return nil, fmt.Errorf("failed to initialize authorization service: %v", err)
		}
	}

	applicationService, err := service.NewApplicationService(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize applications service: %v", err)
	}

	projectsService, err := service.NewProjectsService(
		cfg.Mlflow.TrackingURL,
		repository.NewProjectRepository(db),
		authEnforcer,
		cfg.Authorization.Enabled)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize projects service: %v", err)
	}

	secretService := service.NewSecretService(repository.NewSecretRepository(db, cfg.EncryptionKey))

	return &AppContext{
		ApplicationService:   applicationService,
		ProjectsService:      projectsService,
		SecretService:        secretService,
		AuthorizationEnabled: cfg.Authorization.Enabled,
		Enforcer:             authEnforcer,
	}, nil
}

// type Handler func(r *http.Request, vars map[string]string, body interface{}) *Response
type Handler func(r *http.Request, vars map[string]string, body interface{}) *Response

type Route struct {
	Method  string
	Path    string
	Body    interface{}
	Handler Handler
	Name    string
}

func (route Route) HandlerFunc(validate *validator.Validate) http.HandlerFunc {
	var bodyType reflect.Type
	if route.Body != nil {
		bodyType = reflect.TypeOf(route.Body)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		for k, v := range r.URL.Query() {
			if len(v) > 0 {
				vars[k] = v[0]
			}
		}

		response := func() *Response {
			vars["user"] = r.Header.Get("User-Email")
			var body interface{}

			if bodyType != nil {
				body = reflect.New(bodyType).Interface()
				if err := json.NewDecoder(r.Body).Decode(body); err != nil {
					return BadRequest(fmt.Sprintf("Failed to deserialize request body: %s", err.Error()))
				} else if err := validate.Struct(body); err != nil {
					errMessage := err.(validator.ValidationErrors)[0].Translate(validation.EN)
					return BadRequest(errMessage)
				}
			}
			return route.Handler(r, vars, body)
		}()

		response.WriteTo(w)
	}
}

func NewRouter(appCtx *AppContext, controllers []Controller) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	validator := validation.NewValidator()

	if appCtx.AuthorizationEnabled {
		authzMiddleware := middleware.NewAuthorizer(appCtx.Enforcer)
		router.Use(authzMiddleware.AuthorizationMiddleware)
	}

	for _, c := range controllers {
		for _, r := range c.Routes() {
			_, handler := newrelic.WrapHandle(r.Name, r.HandlerFunc(validator))

			if r.Name == "CreateProject" {
				handler = middleware.ProjectCreationMiddleware(handler)
			}

			router.Name(r.Name).
				Methods(r.Method).
				Path(r.Path).
				Handler(handler)
		}
	}

	return router
}
