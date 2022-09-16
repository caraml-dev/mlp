package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gojek/mlp/api/pkg/authz/enforcer"
	"github.com/gorilla/mux"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/heptiolabs/healthcheck"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"

	"github.com/gojek/mlp/api/api"
	"github.com/gojek/mlp/api/config"
	"github.com/gojek/mlp/api/log"
	"github.com/gojek/mlp/api/service"
	"github.com/gojek/mlp/api/storage"
)

func main() {
	cfg, err := config.InitConfigEnv()
	if err != nil {
		log.Panicf("Failed initializing config: %v", err)
	}

	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			cfg.DbConfig.Host,
			cfg.DbConfig.Port,
			cfg.DbConfig.User,
			cfg.DbConfig.Database,
			cfg.DbConfig.Password))
	if err != nil {
		panic(err)
	}
	db.LogMode(false)
	defer db.Close()

	runDBMigration(db, cfg.DbConfig.MigrationPath)

	applicationService, _ := service.NewApplicationService(db)
	authEnforcer, _ := enforcer.NewEnforcerBuilder().
		URL(cfg.AuthorizationConfig.AuthorizationServerUrl).
		Product("mlp").
		Build()

	projectsService, err := service.NewProjectsService(cfg.MlflowConfig.TrackingUrl, storage.NewProjectStorage(db), authEnforcer, cfg.AuthorizationConfig.AuthorizationEnabled)
	if err != nil {
		log.Panicf("unable to initialize project service: %v", err)
	}

	secretService := service.NewSecretService(storage.NewSecretStorage(db, cfg.EncryptionKey))

	appCtx := api.AppContext{
		ApplicationService: applicationService,
		ProjectsService:    projectsService,
		SecretService:      secretService,

		AuthorizationEnabled: cfg.AuthorizationConfig.AuthorizationEnabled,
		Enforcer:             authEnforcer,
	}

	router := mux.NewRouter()
	mount(router, "/v1/internal", healthcheck.NewHandler())
	mount(router, "/v1", api.NewRouter(appCtx))

	uiEnv := uiEnvHandler{
		ApiURL:        cfg.APIHost,
		OauthClientID: cfg.OauthClientID,
		Environment:   cfg.Environment,
		SentryDSN:     cfg.SentryDSN,
		Teams:         cfg.Teams,
		Streams:       cfg.Streams,
		Docs:          cfg.Docs,

		UIConfig: cfg.UI,
	}

	router.Methods("GET").Path("/env.js").HandlerFunc(uiEnv.handler)

	ui := uiHandler{staticPath: cfg.UI.StaticPath, indexPath: cfg.UI.IndexPath}
	router.PathPrefix("/").Handler(ui)

	log.Infof("listening at port %d", cfg.Port)
	_ = http.ListenAndServe(cfg.ListenAddress(), cors.AllowAll().Handler(router))
}

func mount(r *mux.Router, path string, handler http.Handler) {
	r.PathPrefix(path).Handler(
		http.StripPrefix(
			strings.TrimSuffix(path, "/"),
			handler,
		),
	)
}

type uiEnvHandler struct {
	ApiURL        string                `json:"REACT_APP_API_URL,omitempty"`
	OauthClientID string                `json:"REACT_APP_OAUTH_CLIENT_ID,omitempty"`
	Environment   string                `json:"REACT_APP_ENVIRONMENT,omitempty"`
	SentryDSN     string                `json:"REACT_APP_SENTRY_DSN,omitempty"`
	Teams         []string              `json:"REACT_APP_TEAMS"`
	Streams       []string              `json:"REACT_APP_STREAMS"`
	Docs          config.Documentations `json:"REACT_APP_DOC_LINKS"`

	config.UIConfig
}

func (h uiEnvHandler) handler(w http.ResponseWriter, r *http.Request) {
	envJSON, err := json.Marshal(h)
	if err != nil {
		envJSON = []byte("{}")
	}
	fmt.Fprintf(w, "window.env = %s;", envJSON)
}

// uiHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type uiHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h uiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func runDBMigration(db *gorm.DB, migrationPath string) {
	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
