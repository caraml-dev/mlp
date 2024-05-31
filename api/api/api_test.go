package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/it/database"
	"github.com/caraml-dev/mlp/api/models"
)

type APITestSuite struct {
	suite.Suite

	// handler under test
	route http.Handler

	// db cleanup function
	cleanupFn func()

	// main project to be used for testing
	// most test entities (secret and secret_storage) are created under this
	mainProject *models.Project
	// other project to be used for testing
	otherProject *models.Project

	// interal secret storage reference
	internalSecretStorage *models.SecretStorage
	// default secret storage, the test configures it to vault
	defaultSecretStorage *models.SecretStorage
	// secret storage owned by the project
	projectSecretStorage *models.SecretStorage
	// existing secrets owned by the project
	// in total there are 5 secrets
	// the first 2 secrets are stored in internalSecretStorage
	// the rest of secrets are stored in defaultSecretStorage
	existingSecrets []*models.Secret
}

func (s *APITestSuite) SetupTest() {
	db, cleanupFn, err := database.CreateTestDatabase()
	s.Require().NoError(err, "Failed to connect to test database")
	s.cleanupFn = cleanupFn

	// create app context which also specify vault as the default secret storage
	appCtx, err := NewAppContext(db, &config.Config{
		Port: 0,
		Authorization: &config.AuthorizationConfig{
			Enabled: false,
		},
		Mlflow: &config.MlflowConfig{
			TrackingURL: "http://mlflow:5000",
		},
		UpdateProject: &config.UpdateProjectConfig{
			Endpoint:         "",
			PayloadTemplate:  "template-payload",
			ResponseTemplate: "template-response",
		},
		DefaultSecretStorage: &config.SecretStorage{
			Name: "vault",
			Type: string(models.VaultSecretStorageType),
			Config: models.SecretStorageConfig{
				VaultConfig: &models.VaultConfig{
					URL:        "http://localhost:8200",
					Role:       "my-role",
					MountPath:  "secret",
					PathPrefix: fmt.Sprintf("api-test/%d/{{ .Project }}", time.Now().Unix()),
					AuthMethod: models.TokenAuthMethod,
					Token:      "root",
				},
			},
		},
	})
	s.Require().NoError(err, "Failed to create app context")

	// create project and otherProject
	s.mainProject, err = appCtx.ProjectsService.CreateProject(context.Background(), &models.Project{
		Name: "test-project",
	})
	s.Require().NoError(err, "Failed to create project")

	s.otherProject, err = appCtx.ProjectsService.CreateProject(context.Background(), &models.Project{
		Name: "other-project",
	})
	s.Require().NoError(err, "Failed to create other project")

	// get reference to internal secret storage
	// internal secret storage always have ID 1
	s.internalSecretStorage, err = appCtx.SecretStorageService.FindByID(1)
	s.Require().NoError(err, "Failed to find internal secret storage")

	// get reference to default secret storage
	s.defaultSecretStorage = appCtx.DefaultSecretStorage

	// create a secret storage owned by the project
	s.projectSecretStorage, err = appCtx.SecretStorageService.Create(&models.SecretStorage{
		Name:      "project-secret-storage",
		Type:      models.VaultSecretStorageType,
		Scope:     models.ProjectSecretStorageScope,
		ProjectID: &s.mainProject.ID,
		Config: models.SecretStorageConfig{
			VaultConfig: &models.VaultConfig{
				URL:        "http://localhost:8200",
				Role:       "my-role",
				MountPath:  "secret",
				PathPrefix: fmt.Sprintf("project-secret-test/%d", time.Now().Unix()),
				AuthMethod: models.TokenAuthMethod,
				Token:      "root",
			},
		},
	})
	s.Require().NoError(err, "Failed to create project secret storage")

	//
	s.existingSecrets = make([]*models.Secret, 0)
	for i := 0; i < 5; i++ {
		secretStorageID := s.defaultSecretStorage.ID
		// the first 2 secrets will be created in internal secret storage
		if i < 2 {
			secretStorageID = s.internalSecretStorage.ID
		}

		secret, err := appCtx.SecretService.Create(&models.Secret{
			ProjectID:       s.mainProject.ID,
			SecretStorageID: &secretStorageID,
			Name:            fmt.Sprintf("secret-%d", i),
			Data:            fmt.Sprintf("secret-data-%d", i),
		})
		s.Require().NoError(err, "Failed to create secret")
		s.existingSecrets = append(s.existingSecrets, secret)
	}

	// initialize controllers and http handler / route

	controllers := []Controller{
		&ApplicationsController{AppContext: appCtx},
		&ProjectsController{AppContext: appCtx},
		&SecretsController{AppContext: appCtx},
		&SecretStoragesController{AppContext: appCtx},
	}

	r := NewRouter(appCtx, controllers)

	route := mux.NewRouter()
	route.PathPrefix(basePath).Handler(
		http.StripPrefix(
			strings.TrimSuffix(basePath, "/"),
			r,
		),
	)

	s.route = route
}

func (s *APITestSuite) TearDownTest() {
	s.cleanupFn()
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
