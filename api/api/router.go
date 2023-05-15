package api

import (
	"encoding/json"
	"fmt"
	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/middleware"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
	"github.com/caraml-dev/mlp/api/pkg/instrumentation/newrelic"
	"github.com/caraml-dev/mlp/api/pkg/secretstorage"
	"github.com/caraml-dev/mlp/api/repository"
	"github.com/caraml-dev/mlp/api/service"
	"github.com/caraml-dev/mlp/api/validation"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"net/http"
	"reflect"
)

type Controller interface {
	Routes() []Route
}

type AppContext struct {
	ApplicationService   service.ApplicationService
	ProjectsService      service.ProjectsService
	SecretService        service.SecretService
	SecretStorageService service.SecretStorageService
	DefaultSecretStorage *models.SecretStorage

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

	secretRepository := repository.NewSecretRepository(db)
	storageRepository := repository.NewSecretStorageRepository(db)
	projectRepository := repository.NewProjectRepository(db)

	// initialize default secret storage or create one
	defaultSecretStorage, err := initializeDefaultSecretStorage(storageRepository, cfg)
	if err != nil {
		return nil, err
	}

	// get all secret storages and create corresponding clients
	allSecretStorages, err := storageRepository.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list all secret storages: %v", err)
	}
	storageClientRegistry, err := secretstorage.NewRegistry(allSecretStorages)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize secret storage registry: %v", err)
	}

	secretService := service.NewSecretService(secretRepository, storageRepository, projectRepository, storageClientRegistry, defaultSecretStorage)
	secretStorageService := service.NewSecretStorageService(storageRepository, storageClientRegistry)

	return &AppContext{
		ApplicationService:   applicationService,
		ProjectsService:      projectsService,
		SecretService:        secretService,
		SecretStorageService: secretStorageService,
		AuthorizationEnabled: cfg.Authorization.Enabled,
		Enforcer:             authEnforcer,
		DefaultSecretStorage: defaultSecretStorage,
	}, nil
}

func initializeDefaultSecretStorage(storageRepository repository.SecretStorageRepository, cfg *config.Config) (*models.SecretStorage, error) {
	defaultSecretStorage, err := storageRepository.GetGlobal(cfg.DefaultSecretStorageModel().Name)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to initialize default secret storage: %v", err)
		}

		// create one if not found
		return storageRepository.Save(cfg.DefaultSecretStorageModel())
	}

	// update default secret storage if it has changed
	err = copier.CopyWithOption(defaultSecretStorage, cfg.DefaultSecretStorageModel(), copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return storageRepository.Save(defaultSecretStorage)
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
