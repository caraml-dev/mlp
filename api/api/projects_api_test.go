//go:build integration

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mux2 "github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/caraml-dev/mlp/api/config"
	"github.com/caraml-dev/mlp/api/it/database"
	"github.com/caraml-dev/mlp/api/models"
	"github.com/caraml-dev/mlp/api/repository"
	"github.com/caraml-dev/mlp/api/service"
)

const (
	mlflowTrackingURL = "http://localhost.com"
	adminUser         = "admin"
	basePath          = "/v1"
)

var (
	now = time.Now()
)

func TestCreateProject(t *testing.T) {
	testCases := []struct {
		desc             string
		userAgent        string
		existingProject  *models.Project
		body             interface{}
		errSaveSecret    error
		expectedResponse *Response
	}{
		{
			desc:      "Should success for project without labels",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.97 Safari/537.36]",
			body: &models.Project{
				Name:   "my-project",
				Team:   "dsp",
				Stream: "dsp",
			},
			expectedResponse: &Response{
				code: 201,
				data: &models.Project{
					ID:                models.ID(1),
					Name:              "my-project",
					MLFlowTrackingURL: mlflowTrackingURL,
					Administrators:    []string{adminUser},
					Team:              "dsp",
					Stream:            "dsp",
					CreatedUpdated: models.CreatedUpdated{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			existingProject: nil,
		},
		{
			desc:      "Should success for project with labels",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.97 Safari/537.36]",
			body: &models.Project{
				Name:   "my-project",
				Team:   "dsp",
				Stream: "dsp",
				Labels: models.Labels{
					{
						Key:   "my-label",
						Value: "my-value",
					},
				},
			},
			expectedResponse: &Response{
				code: 201,
				data: &models.Project{
					ID:                models.ID(1),
					Name:              "my-project",
					MLFlowTrackingURL: mlflowTrackingURL,
					Administrators:    []string{adminUser},
					Team:              "dsp",
					Stream:            "dsp",
					Labels: models.Labels{
						{
							Key:   "my-label",
							Value: "my-value",
						},
					},
					CreatedUpdated: models.CreatedUpdated{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			existingProject: nil,
		},
		{
			desc:      "Should fail when project with same name exists",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.97 Safari/537.36]",
			body: &models.Project{
				Name:   "my-project",
				Team:   "dsp",
				Stream: "dsp",
			},
			expectedResponse: &Response{
				code: 400,
				data: ErrorMessage{"Project my-project already exists"},
			},
			existingProject: &models.Project{
				Name:   "my-project",
				Team:   "dsp",
				Stream: "dsp",
			},
		},
		{
			desc:      "Should fail when request doesn't specify team",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.97 Safari/537.36]",
			body: &models.Project{
				Name:   "my-project",
				Stream: "dsp",
			},
			expectedResponse: &Response{
				code: 400,
				data: ErrorMessage{"Team is required"},
			},
			existingProject: nil,
		},
		{
			desc:      "Should fail when request doesn't specify stream",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.97 Safari/537.36]",
			body: &models.Project{
				Name: "my-project",
				Team: "dsp",
			},
			expectedResponse: &Response{
				code: 400,
				data: ErrorMessage{"Stream is required"},
			},
			existingProject: nil,
		},
		{
			desc:      "Should fail when project name is shorter than 3 characters",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.97 Safari/537.36]",
			body: &models.Project{
				Name: "a",
				Team: "dsp",
			},
			expectedResponse: &Response{
				code: 400,
				data: ErrorMessage{"Name should be more than 3 characters"},
			},
			existingProject: nil,
		},
		{
			desc:      "Should fail when project name is longer than 50 characters",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.97 Safari/537.36]",
			body: &models.Project{
				Name: "lorem-ipsum-dolor-sing-amet-hahaha-hihihi-huhuhu-hehe-hoho",
				Team: "dsp",
			},
			expectedResponse: &Response{
				code: 400,
				data: ErrorMessage{"Name should be less than 50 characters"},
			},
			existingProject: nil,
		},
		{
			desc:      "Should fail when project name is not RFC1123 compliant",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.97 Safari/537.36]",
			body: &models.Project{
				Name: "-invalid-project",
				Team: "dsp",
			},
			expectedResponse: &Response{
				code: 400,
				data: ErrorMessage{"Name should be a valid RFC1123 sub-domain"},
			},
			existingProject: nil,
		},
		{
			desc:      "Should fail if user agent is swagger-codegen",
			userAgent: "Swagger-Codegen/1.0.0/python",
			body: &models.Project{
				Name: "swag-project",
			},
			expectedResponse: &Response{
				code: 403,
				data: ErrorMessage{"Project creation from SDK is disabled. Use the MLP console to create a project."},
			},
			existingProject: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				prjRepository := repository.NewProjectRepository(db)
				if tC.existingProject != nil {
					_, err := prjRepository.Save(tC.existingProject)
					assert.NoError(t, err)
				}
				projectService, err := service.NewProjectsService(
					mlflowTrackingURL, prjRepository, nil, false, nil,
					config.UpdateProjectConfig{},
				)
				assert.NoError(t, err)

				appCtx := &AppContext{
					ProjectsService:      projectService,
					AuthorizationEnabled: false,
				}
				controllers := []Controller{&ProjectsController{appCtx}}
				r := NewRouter(appCtx, controllers)

				requestByte, _ := json.Marshal(tC.body)
				req, err := http.NewRequest(http.MethodPost, "/v1/projects", bytes.NewReader(requestByte))
				if err != nil {
					t.Fatal(err)
				}

				req.Header["User-Email"] = []string{adminUser}
				req.Header["User-Agent"] = []string{tC.userAgent}
				rr := httptest.NewRecorder()

				route := mux2.NewRouter()
				route.PathPrefix(basePath).Handler(
					http.StripPrefix(
						strings.TrimSuffix(basePath, "/"),
						r,
					),
				)
				route.ServeHTTP(rr, req)

				assert.Equal(t, tC.expectedResponse.code, rr.Code)
				if tC.expectedResponse.code >= 200 && tC.expectedResponse.code < 300 {
					project := &models.Project{}
					err = json.Unmarshal(rr.Body.Bytes(), &project)
					assert.NoError(t, err)

					project.CreatedAt = now
					project.UpdatedAt = now

					assert.Equal(t, tC.expectedResponse.data, project)
				} else {
					e := ErrorMessage{}
					err = json.Unmarshal(rr.Body.Bytes(), &e)
					assert.NoError(t, err)

					assert.Equal(t, tC.expectedResponse.data, e)
				}
			})
		})
	}
}

func TestListProjects(t *testing.T) {
	testCases := []struct {
		desc             string
		existingProjects []models.Project
		expectedResponse *Response
	}{
		{
			desc: "Should return all",
			existingProjects: []models.Project{
				{
					ID:                models.ID(1),
					Name:              "Project1",
					MLFlowTrackingURL: "http://mlflow.com",
					Administrators:    []string{adminUser},
					Team:              "dsp",
					Stream:            "dsp",
					CreatedUpdated: models.CreatedUpdated{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			expectedResponse: &Response{
				code: 200,
				data: []*models.Project{
					{
						ID:                models.ID(1),
						Name:              "Project1",
						MLFlowTrackingURL: "http://mlflow.com",
						Administrators:    []string{adminUser},
						Team:              "dsp",
						Stream:            "dsp",
						CreatedUpdated: models.CreatedUpdated{
							CreatedAt: now,
							UpdatedAt: now,
						},
					},
				},
			},
		},
		{
			desc:             "Should return empty project",
			existingProjects: []models.Project{},
			expectedResponse: &Response{
				code: 200,
				data: []*models.Project{},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				prjRepository := repository.NewProjectRepository(db)
				if tC.existingProjects != nil {
					for _, project := range tC.existingProjects {
						_, err := prjRepository.Save(&project)
						assert.NoError(t, err)
					}
				}
				projectService, err := service.NewProjectsService(
					mlflowTrackingURL, prjRepository, nil, false, nil,
					config.UpdateProjectConfig{},
				)
				assert.NoError(t, err)

				appCtx := &AppContext{
					ProjectsService:      projectService,
					AuthorizationEnabled: false,
				}
				controllers := []Controller{&ProjectsController{appCtx}}
				r := NewRouter(appCtx, controllers)

				req, err := http.NewRequest(http.MethodGet, "/v1/projects", nil)
				if err != nil {
					t.Fatal(err)
				}

				req.Header["User-Email"] = []string{adminUser}
				rr := httptest.NewRecorder()

				route := mux2.NewRouter()
				route.PathPrefix(basePath).Handler(
					http.StripPrefix(
						strings.TrimSuffix(basePath, "/"),
						r,
					),
				)
				route.ServeHTTP(rr, req)

				assert.Equal(t, tC.expectedResponse.code, rr.Code)
				if tC.expectedResponse.code >= 200 && tC.expectedResponse.code < 300 {
					projects := []*models.Project{}
					err = json.Unmarshal(rr.Body.Bytes(), &projects)
					assert.NoError(t, err)

					for _, project := range projects {
						project.CreatedAt = now
						project.UpdatedAt = now
					}

					assert.Equal(t, tC.expectedResponse.data, projects)
				} else {
					e := ErrorMessage{}
					err = json.Unmarshal(rr.Body.Bytes(), &e)
					assert.NoError(t, err)

					assert.Equal(t, tC.expectedResponse.data, e)
				}
			})
		})
	}
}

func TestUpdateProject(t *testing.T) {
	testCases := []struct {
		desc                string
		projectID           models.ID
		existingProject     *models.Project
		expectedResponse    *Response
		body                interface{}
		updateProjectConfig config.UpdateProjectConfig
	}{
		{
			desc:      "Should success with update project config",
			projectID: models.ID(1),
			existingProject: &models.Project{
				ID:                models.ID(1),
				Name:              "Project1",
				MLFlowTrackingURL: "http://mlflow.com",
				Administrators:    []string{adminUser},
				Team:              "dsp",
				Stream:            "dsp",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			body: &models.Project{
				Name:           "Project1",
				Team:           "merlin",
				Stream:         "dsp",
				Administrators: []string{adminUser},
			},
			expectedResponse: &Response{
				code: 200,
				data: map[string]interface{}{
					"status":  "success",
					"message": "Project updated successfully",
				},
			},
			updateProjectConfig: config.UpdateProjectConfig{
				Endpoint: "url",
				PayloadTemplate: `{
					"project": "{{.Name}}",
					"administrators": "{{.Administrators}}",
					"readers": "{{.Readers}}",
					"team": "{{.Team}}",
					"stream": "{{.Stream}}"
				}`,
				ResponseTemplate: `{
					"status": "{{.status}}",
					"message": "{{.message}}"
				}`,
				LabelsBlacklist: []string{
					"label1",
					"label2",
				},
			},
		},
		{
			desc:      "Should success without update project config",
			projectID: models.ID(1),
			existingProject: &models.Project{
				ID:                models.ID(1),
				Name:              "Project1",
				MLFlowTrackingURL: "http://mlflow.com",
				Administrators:    []string{adminUser},
				Team:              "dsp",
				Stream:            "dsp",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			body: &models.Project{
				Name:           "Project1",
				Team:           "merlin",
				Stream:         "dsp",
				Administrators: []string{adminUser},
			},
			expectedResponse: &Response{
				code: 200,
				data: &models.Project{
					ID:                models.ID(1),
					Name:              "Project1",
					MLFlowTrackingURL: "http://mlflow.com",
					Administrators:    []string{adminUser},
					Team:              "merlin",
					Stream:            "dsp",
					CreatedUpdated: models.CreatedUpdated{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			updateProjectConfig: config.UpdateProjectConfig{},
		},
		{
			desc:      "Should failed when name is not specified",
			projectID: models.ID(1),
			existingProject: &models.Project{
				ID:                models.ID(1),
				Name:              "Project1",
				MLFlowTrackingURL: "http://mlflow.com",
				Administrators:    []string{adminUser},
				Team:              "dsp",
				Stream:            "dsp",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			body: &models.Project{
				Team:           "merlin",
				Stream:         "dsp",
				Administrators: []string{adminUser},
			},
			expectedResponse: &Response{
				code: 400,
				data: ErrorMessage{
					Message: "Name is required",
				},
			},
			updateProjectConfig: config.UpdateProjectConfig{},
		},
		{
			desc:      "Should failed when name project id is not found",
			projectID: models.ID(2),
			existingProject: &models.Project{
				ID:                models.ID(1),
				Name:              "Project1",
				MLFlowTrackingURL: "http://mlflow.com",
				Administrators:    []string{adminUser},
				Team:              "dsp",
				Stream:            "dsp",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			body: &models.Project{
				Name:           "Project1",
				Team:           "merlin",
				Stream:         "dsp",
				Administrators: []string{adminUser},
			},
			expectedResponse: &Response{
				code: 404,
				data: ErrorMessage{
					Message: "project with ID 2 not found",
				},
			},
			updateProjectConfig: config.UpdateProjectConfig{},
		},
		{
			desc:      "Should fail when label in blacklist",
			projectID: models.ID(1),
			existingProject: &models.Project{
				ID:                models.ID(1),
				Name:              "Project1",
				MLFlowTrackingURL: "http://mlflow.com",
				Administrators:    []string{adminUser},
				Team:              "dsp",
				Stream:            "dsp",
				Labels: models.Labels{
					{
						Key:   "my-label",
						Value: "my-value",
					},
				},
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			body: &models.Project{
				Name:   "Project1",
				Team:   "merlin",
				Stream: "dsp",
				Labels: models.Labels{
					{
						Key:   "my-label",
						Value: "my-new-value",
					},
				},
				Administrators: []string{adminUser},
			},
			expectedResponse: &Response{
				code: 500,
				data: ErrorMessage{
					Message: "one or more labels are blacklisted or have been removed or changed values and cannot be updated",
				},
			},
			updateProjectConfig: config.UpdateProjectConfig{
				LabelsBlacklist: []string{
					"my-label",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				prjRepository := repository.NewProjectRepository(db)
				if tC.existingProject != nil {
					_, err := prjRepository.Save(tC.existingProject)
					assert.NoError(t, err)
				}

				var server *httptest.Server
				if tC.updateProjectConfig.Endpoint != "" {
					server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var payload map[string]interface{}
						err := json.NewDecoder(r.Body).Decode(&payload)
						assert.NoError(t, err)

						w.WriteHeader(http.StatusOK)
						response := map[string]string{
							"status":  "success",
							"message": "Project updated successfully",
						}
						err = json.NewEncoder(w).Encode(response)
						assert.NoError(t, err)
					}))
					defer server.Close()

					tC.updateProjectConfig.Endpoint = server.URL
				}

				projectService, err := service.NewProjectsService(
					mlflowTrackingURL, prjRepository, nil, false, nil,
					tC.updateProjectConfig,
				)
				assert.NoError(t, err)

				appCtx := &AppContext{
					ProjectsService:      projectService,
					AuthorizationEnabled: false,
				}
				controllers := []Controller{&ProjectsController{appCtx}}
				r := NewRouter(appCtx, controllers)

				requestByte, _ := json.Marshal(tC.body)
				req, err := http.NewRequest(http.MethodPut, "/v1/projects/"+tC.projectID.String(), bytes.NewReader(requestByte))
				if err != nil {
					t.Fatal(err)
				}

				req.Header["User-Email"] = []string{adminUser}
				rr := httptest.NewRecorder()

				route := mux2.NewRouter()
				route.PathPrefix(basePath).Handler(
					http.StripPrefix(
						strings.TrimSuffix(basePath, "/"),
						r,
					),
				)
				route.ServeHTTP(rr, req)

				assert.Equal(t, tC.expectedResponse.code, rr.Code)
				if tC.expectedResponse.code >= 200 && tC.expectedResponse.code < 300 {
					switch tC.expectedResponse.data.(type) {
					case *models.Project:
						project := &models.Project{}
						err = json.Unmarshal(rr.Body.Bytes(), &project)
						assert.NoError(t, err)

						project.CreatedAt = now
						project.UpdatedAt = now

						assert.Equal(t, tC.expectedResponse.data, project)
					case map[string]interface{}:
						var responseMessage map[string]interface{}
						err = json.Unmarshal(rr.Body.Bytes(), &responseMessage)
						assert.NoError(t, err)
						assert.Equal(t, tC.expectedResponse.data, responseMessage)
					default:
						t.Fatal("unexpected type for expectedResponse.data")
					}
				} else {
					e := ErrorMessage{}
					err = json.Unmarshal(rr.Body.Bytes(), &e)
					assert.NoError(t, err)

					assert.Equal(t, tC.expectedResponse.data, e)
				}
			})
		})
	}
}

func TestGetProject(t *testing.T) {
	testCases := []struct {
		desc             string
		projectID        models.ID
		existingProject  *models.Project
		expectedResponse *Response
	}{
		{
			desc:      "Should success",
			projectID: models.ID(1),
			existingProject: &models.Project{
				ID:                models.ID(1),
				Name:              "Project1",
				MLFlowTrackingURL: "http://mlflow.com",
				Administrators:    []string{adminUser},
				Team:              "dsp",
				Stream:            "dsp",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expectedResponse: &Response{
				code: 200,
				data: &models.Project{
					ID:                models.ID(1),
					Name:              "Project1",
					MLFlowTrackingURL: "http://mlflow.com",
					Administrators:    []string{adminUser},
					Team:              "dsp",
					Stream:            "dsp",
					CreatedUpdated: models.CreatedUpdated{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
		},
		{
			desc:      "Should return nothing if project is not exist",
			projectID: models.ID(2),
			existingProject: &models.Project{
				ID:                models.ID(1),
				Name:              "Project1",
				MLFlowTrackingURL: "http://mlflow.com",
				Administrators:    []string{adminUser},
				Team:              "dsp",
				Stream:            "dsp",
				CreatedUpdated: models.CreatedUpdated{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expectedResponse: &Response{
				code: 404,
				data: ErrorMessage{
					Message: "project with ID 2 not found",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			database.WithTestDatabase(t, func(t *testing.T, db *gorm.DB) {
				prjRepository := repository.NewProjectRepository(db)
				if tC.existingProject != nil {
					_, err := prjRepository.Save(tC.existingProject)
					assert.NoError(t, err)
				}
				projectService, err := service.NewProjectsService(
					mlflowTrackingURL, prjRepository, nil, false, nil,
					config.UpdateProjectConfig{},
				)
				assert.NoError(t, err)

				appCtx := &AppContext{
					ProjectsService:      projectService,
					AuthorizationEnabled: false,
				}
				controllers := []Controller{&ProjectsController{appCtx}}
				r := NewRouter(appCtx, controllers)

				req, err := http.NewRequest(http.MethodGet, "/v1/projects/"+tC.projectID.String(), nil)
				if err != nil {
					t.Fatal(err)
				}

				req.Header["User-Email"] = []string{adminUser}
				rr := httptest.NewRecorder()

				route := mux2.NewRouter()
				route.PathPrefix(basePath).Handler(
					http.StripPrefix(
						strings.TrimSuffix(basePath, "/"),
						r,
					),
				)
				route.ServeHTTP(rr, req)

				assert.Equal(t, tC.expectedResponse.code, rr.Code)
				if tC.expectedResponse.code >= 200 && tC.expectedResponse.code < 300 {
					project := &models.Project{}
					err = json.Unmarshal(rr.Body.Bytes(), &project)
					assert.NoError(t, err)

					project.CreatedAt = now
					project.UpdatedAt = now

					assert.Equal(t, tC.expectedResponse.data, project)
				} else {
					e := ErrorMessage{}
					err = json.Unmarshal(rr.Body.Bytes(), &e)
					assert.NoError(t, err)

					assert.Equal(t, tC.expectedResponse.data, e)
				}
			})
		})
	}
}
