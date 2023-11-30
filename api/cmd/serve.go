package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gojekfarm/iap_auth/pkg/iap"
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/caraml-dev/mlp/api/api"
	apiV2 "github.com/caraml-dev/mlp/api/api/v2"
	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/database"
	"github.com/caraml-dev/mlp/api/log"
	"github.com/caraml-dev/mlp/api/pkg/authz/enforcer"
)

var (
	configFiles []string
	serveCmd    = &cobra.Command{
		Use:   "serve",
		Short: "Start MLP API server",
		Run: func(cmd *cobra.Command, args []string) {
			serveConfig, err := config.LoadAndValidate(configFiles...)
			if err != nil {
				log.Fatalf("failed initializing config: %v", err)
			}
			startServer(serveConfig)
		},
	}
)

func init() {
	serveCmd.Flags().StringSliceVarP(&configFiles, "config", "c", []string{},
		"Comma separated list of config files to load. The last config file will take precedence over the "+
			"previous ones.")
}

func startServer(cfg *config.Config) {
	log.Infof("starting MLP API server with config: %+v", cfg.Alert)

	// init db
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Panicf("unable to initialize DB connectivity: %v", err)
	}
	defer db.Close()

	hc := oauth2.NewClient(context.Background(), nil)
	tokenfn := func() string {
		iap, err := iap.New(hc, cfg.Alert.ServiceAccountCredentials, cfg.Alert.ClientID)
		if err != nil {
			return "INVALID"
		}
		token, err := iap.Token()
		if err != nil {
			return "INVALID"
		}
		return token
	}

	var atomictoken atomic.Value
	atomictoken.Store(tokenfn())

	ticker := time.NewTicker(cfg.Alert.TickPeriod)
	var mu sync.Mutex
	go func() {
		for range ticker.C {
			mu.Lock()
			atomictoken.Store(tokenfn())
			mu.Unlock()
		}
	}()

	appCtx, err := api.NewAppContext(&atomictoken, db, cfg)
	if err != nil {
		log.Panicf("unable to initialize application context: %v", err)
	}

	router := mux.NewRouter()

	mount(router, "/v1/internal", healthcheck.NewHandler())

	v1Controllers := []api.Controller{
		&api.AlertController{AppContext: appCtx},
		&api.ApplicationsController{AppContext: appCtx},
		&api.ProjectsController{AppContext: appCtx},
		&api.SecretsController{AppContext: appCtx},
		&api.SecretStoragesController{AppContext: appCtx},
	}
	mount(router, "/v1", api.NewRouter(appCtx, v1Controllers))

	v2Controllers := []api.Controller{
		&apiV2.ApplicationsController{Apps: cfg.Applications},
	}
	mount(router, "/v2", api.NewRouter(appCtx, v2Controllers))

	uiEnv := uiEnvHandler{
		APIURL:        cfg.APIHost,
		OauthClientID: cfg.OauthClientID,
		Environment:   cfg.Environment,
		SentryDSN:     cfg.SentryDSN,
		Streams:       cfg.Streams,
		Docs:          cfg.Docs,
		MaxAuthzCacheExpiryMinutes: fmt.Sprintf("%.0f",
			math.Ceil((time.Duration(enforcer.MaxKeyExpirySeconds) * time.Second).Minutes())),
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
	*config.UIConfig

	APIURL                     string                `json:"REACT_APP_API_URL,omitempty"`
	OauthClientID              string                `json:"REACT_APP_OAUTH_CLIENT_ID,omitempty"`
	Environment                string                `json:"REACT_APP_ENVIRONMENT,omitempty"`
	SentryDSN                  string                `json:"REACT_APP_SENTRY_DSN,omitempty"`
	Streams                    config.Streams        `json:"REACT_APP_STREAMS"`
	Docs                       config.Documentations `json:"REACT_APP_DOC_LINKS"`
	MaxAuthzCacheExpiryMinutes string                `json:"REACT_APP_MAX_AUTHZ_CACHE_EXPIRY_MINUTES"`
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
