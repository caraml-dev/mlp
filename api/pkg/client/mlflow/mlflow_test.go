package mlflow

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/storage"

	"github.com/caraml-dev/mlp/api/pkg/artifact"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/caraml-dev/mlp/api/pkg/artifact/mocks"
)

var MultipleRunSuccessJSON = `{    
	"runs": [
		{
			"info": {
				"run_uuid": "run-123",
				"experiment_id": "1",
				"user_id": "root",
				"status": "FINISHED",
				"start_time": "1677735900543",
				"end_time": "1677735901790",
				"artifact_uri": "gs://my-bucket/run-123",
				"lifecycle_stage": "active",
				"run_id": "run-123"
			},
			"data": {
				"tags": [
					{
						"key": "env",
						"value": "prod"
					},
					{
						"key": "version",
						"value": "1.0.0"
					}
				]
			}
		},
		{
			"info": {
				"run_uuid": "run-456",
				"experiment_id": "1",
				"user_id": "root",
				"status": "FINISHED",
				"start_time": "1677735900543",
				"end_time": "1677735901790",
				"artifact_uri": "gs://my-bucket/run-456",
				"lifecycle_stage": "active",
				"run_id": "run-456"
			},
			"data": {
				"tags": [
					{
						"key": "env",
						"value": "dev"
					},
					{
						"key": "version",
						"value": "1.1.0"
					}
				]
			}
		}
	]}`
var MultipleRunSuccessJSONFailedDelete = `{    
	"runs": [
		{
			"info": {
				"run_uuid": "run-123",
				"experiment_id": "1",
				"user_id": "root",
				"status": "FINISHED",
				"start_time": "1677735900543",
				"end_time": "1677735901790",
				"artifact_uri": "gs://my-bucket/run-789",
				"lifecycle_stage": "active",
				"run_id": "run-123"
			},
			"data": {
				"tags": [
					{
						"key": "env",
						"value": "prod"
					},
					{
						"key": "version",
						"value": "1.0.0"
					}
				]
			}
		},
		{
			"info": {
				"run_uuid": "run-456",
				"experiment_id": "1",
				"user_id": "root",
				"status": "FINISHED",
				"start_time": "1677735900543",
				"end_time": "1677735901790",
				"artifact_uri": "gs://my-bucket/run-123",
				"lifecycle_stage": "active",
				"run_id": "run-456"
			},
			"data": {
				"tags": [
					{
						"key": "env",
						"value": "dev"
					},
					{
						"key": "version",
						"value": "1.1.0"
					}
				]
			}
		}
	]}`

var RunSuccessJSON = `
{
	"run": {
		"info": {
			"run_uuid": "run-123",
			"experiment_id": "1",
			"user_id": "root",
			"status": "FINISHED",
			"start_time": "1677735900543",
			"end_time": "1677735901790",
			"artifact_uri": "gs://my-bucket/run-123",
			"lifecycle_stage": "active",
			"run_id": "run-123"
		},
		"data": {
			"tags": [
				{
					"key": "env",
					"value": "prod"
				},
				{
					"key": "version",
					"value": "1.0.0"
				}
			]
		}		
	}
}`

var RunSuccessDeleteRun = `
{
	"run": {
		"info": {
			"run_uuid": "run-123",
			"experiment_id": "1",
			"user_id": "root",
			"status": "FINISHED",
			"start_time": "1677735900543",
			"end_time": "1677735901790",
			"artifact_uri": "gs://bucketName/valid",
			"lifecycle_stage": "active",
			"run_id": "run-123"
		},
		"data": {
			"tags": [
				{
					"key": "env",
					"value": "prod"
				},
				{
					"key": "version",
					"value": "1.0.0"
				}
			]
		}		
	}
}`

var RunFailedDeleteRun = `
{
	"run": {
		"info": {
			"run_uuid": "run-123",
			"experiment_id": "1",
			"user_id": "root",
			"status": "FINISHED",
			"start_time": "1677735900543",
			"end_time": "1677735901790",
			"artifact_uri": "gs://bucketName/invalid",
			"lifecycle_stage": "active",
			"run_id": "run-123"
		},
		"data": {
			"tags": [
				{
					"key": "env",
					"value": "prod"
				},
				{
					"key": "version",
					"value": "1.0.0"
				}
			]
		}		
	}
}`

var RunDoesntExist = `
{
    "error_code": "RESOURCE_DOES_NOT_EXIST",
    "message": "run with id=unknownId not found"
}`

var DeleteRunAlreadyDeleted = `
{
    "error_code": "INVALID_PARAMETER_VALUE",
    "message": "the run xytspow3412oi must be in the 'active' state. Current state is deleted"
}`

func TestNewMlflowService(t *testing.T) {
	httpClient := http.Client{}
	ctx := context.Background()
	api, _ := storage.NewClient(ctx)

	tests := []struct {
		name           string
		artifactType   string
		expectedError  error
		expectedResult *mlflowService
	}{
		{
			name:          "Mlflow Service with GCS Artifact",
			artifactType:  "gcs",
			expectedError: nil,
			expectedResult: &mlflowService{
				API:    &httpClient,
				Config: Config{TrackingURL: "", ArtifactServiceType: "gcs"},
				ArtifactService: &artifact.GcsArtifactClient{
					API: api,
				},
			},
		},
		{
			name:          "Mlflow Service with nop Artifact",
			artifactType:  "nop",
			expectedError: nil,
			expectedResult: &mlflowService{
				API:             &httpClient,
				Config:          Config{TrackingURL: "", ArtifactServiceType: "nop"},
				ArtifactService: &artifact.NopArtifactClient{},
			},
		},
		{
			name:           "Mlflow Service with other Artifact",
			artifactType:   "other",
			expectedError:  fmt.Errorf("invalid artifact service type"),
			expectedResult: &mlflowService{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mlflowService, err := NewMlflowService(&httpClient, Config{
				TrackingURL:         "",
				ArtifactServiceType: tc.artifactType,
			})

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			}
			assert.IsType(t, tc.expectedResult, mlflowService)
		})
	}
}

func TestMlflowClient_SearchRunForExperiment(t *testing.T) {
	tests := []struct {
		name             string
		idExperiment     string
		expectedRespJSON string
		expectedResponse SearchRunsResponse
		expectedError    error
	}{
		{
			name:             "Valid Search",
			idExperiment:     "1",
			expectedRespJSON: MultipleRunSuccessJSON,
			expectedResponse: SearchRunsResponse{
				RunsData: []RunResponse{
					{
						Info: RunInfo{
							RunID:          "run-123",
							ExperimentID:   "1",
							UserID:         "root",
							LifecycleStage: "active",
							ArtifactURI:    "gs://my-bucket/run-123",
						},
						Data: RunData{
							Tags: []RunTag{
								{Key: "env", Value: "prod"},
								{Key: "version", Value: "1.0.0"},
							},
						},
					},
					{
						Info: RunInfo{
							RunID:          "run-456",
							ExperimentID:   "1",
							UserID:         "root",
							LifecycleStage: "active",
							ArtifactURI:    "gs://my-bucket/run-456",
						},
						Data: RunData{
							Tags: []RunTag{
								{Key: "env", Value: "dev"},
								{Key: "version", Value: "1.1.0"},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:             "No related runs",
			idExperiment:     "999",
			expectedRespJSON: `{}`,
			expectedResponse: SearchRunsResponse{},
			expectedError:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(tc.expectedRespJSON))
				require.NoError(t, err)
			}))
			defer server.Close()

			client := mlflowService{
				API:             server.Client(),
				ArtifactService: &mocks.Service{},
				Config:          Config{TrackingURL: server.URL},
			}

			resp, errAPI := client.searchRunsForExperiment(tc.idExperiment)
			assert.Equal(t, tc.expectedError, errAPI)
			assert.Equal(t, tc.expectedResponse, resp)

		})
	}

}

func TestMlflowClient_SearchRunData(t *testing.T) {
	tests := []struct {
		name             string
		idRun            string
		expectedRespJSON string
		expectedResponse SearchRunResponse
		expectedError    error
		httpStatus       int
	}{
		{
			name:             "Valid Search",
			idRun:            "abcdefg1234",
			expectedRespJSON: RunSuccessJSON,
			expectedResponse: SearchRunResponse{
				RunData: RunResponse{
					Info: RunInfo{
						RunID:          "run-123",
						ExperimentID:   "1",
						UserID:         "root",
						LifecycleStage: "active",
						ArtifactURI:    "gs://my-bucket/run-123",
					},
					Data: RunData{
						Tags: []RunTag{
							{Key: "env", Value: "prod"},
							{Key: "version", Value: "1.0.0"},
						},
					},
				},
			},
			expectedError: nil,
			httpStatus:    http.StatusOK,
		},
		{
			name:             "No related runs",
			idRun:            "unknownID",
			expectedRespJSON: RunDoesntExist,
			expectedResponse: SearchRunResponse{},
			expectedError:    fmt.Errorf("run with id=unknownId not found"),
			httpStatus:       http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(tc.httpStatus)
				_, err := w.Write([]byte(tc.expectedRespJSON))
				require.NoError(t, err)
			}))
			defer server.Close()

			client := mlflowService{
				API:             server.Client(),
				ArtifactService: &mocks.Service{},
				Config:          Config{TrackingURL: server.URL},
			}

			resp, errAPI := client.searchRunData(tc.idRun)

			assert.Equal(t, tc.expectedError, errAPI)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}

func TestMlflowClient_DeleteExperiment(t *testing.T) {
	tests := []struct {
		name                 string
		idExperiment         string
		expectedRespJSON     string
		expectedError        error
		httpStatus           int
		expectedRunsRespJSON string
	}{
		{
			name:                 "Valid Experiment Deletion",
			idExperiment:         "1",
			expectedRespJSON:     `{}`,
			expectedError:        nil,
			httpStatus:           http.StatusOK,
			expectedRunsRespJSON: MultipleRunSuccessJSON,
		},
		{
			name:             "Run Failed Deletion",
			idExperiment:     "1",
			expectedRespJSON: `{}`,
			expectedError: fmt.Errorf("deletion failed for run_id run-123 for experiment id 1: " +
				"failed to Delete Artifact"),
			httpStatus:           http.StatusOK,
			expectedRunsRespJSON: MultipleRunSuccessJSONFailedDelete,
		},
		{
			name:                 "No related run for Id",
			idExperiment:         "999",
			expectedRespJSON:     `{}`,
			expectedError:        fmt.Errorf("there are no related run for experiment id 999"),
			httpStatus:           http.StatusOK,
			expectedRunsRespJSON: `{}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/2.0/mlflow/runs/delete", func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(tc.httpStatus)
				_, err := w.Write([]byte(tc.expectedRespJSON))
				require.NoError(t, err)
			})
			mux.HandleFunc("/api/2.0/mlflow/runs/search", func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(tc.httpStatus)
				_, err := w.Write([]byte(tc.expectedRunsRespJSON))
				require.NoError(t, err)
			})

			server := httptest.NewServer(mux)
			defer server.Close()

			artifactServiceMock := mocks.Service{}
			client := mlflowService{
				API:             server.Client(),
				ArtifactService: &artifactServiceMock,
				Config:          Config{TrackingURL: server.URL},
			}

			artifactServiceMock.
				On("DeleteArtifact", context.Background(), "gs://my-bucket/run-789").
				Return(fmt.Errorf("failed to Delete Artifact"))
			artifactServiceMock.
				On("DeleteArtifact", context.Background(), "gs://my-bucket/run-123").
				Return(nil)
			artifactServiceMock.
				On("DeleteArtifact", context.Background(), "gs://my-bucket/run-456").
				Return(nil)

			errAPI := client.DeleteExperiment(context.Background(), tc.idExperiment, true)

			assert.Equal(t, tc.expectedError, errAPI)

		})
	}
}

func TestMlflowClient_DeleteRun(t *testing.T) {
	tests := []struct {
		name                string
		idRun               string
		expectedRespJSON    string
		expectedRunRespJSON string
		expectedError       error
		httpStatus          int
		artifactURL         string
		deleteArtifact      bool
	}{
		{
			name:                "Valid Run Deletion Without Delete Artifact",
			idRun:               "abcdefg1234",
			expectedRespJSON:    `{}`,
			expectedError:       nil,
			httpStatus:          http.StatusOK,
			artifactURL:         "gs://bucketName/valid",
			deleteArtifact:      false,
			expectedRunRespJSON: `{}`,
		},
		{
			name:                "Valid Run Deletion With Delete Artifact",
			idRun:               "abcdefg1234",
			expectedRespJSON:    `{}`,
			expectedError:       nil,
			httpStatus:          http.StatusOK,
			artifactURL:         "gs://bucketName/valid",
			deleteArtifact:      true,
			expectedRunRespJSON: `{}`,
		},
		{
			name:                "ID already deleted",
			idRun:               "xytspow3412oi",
			expectedRespJSON:    DeleteRunAlreadyDeleted,
			expectedError:       fmt.Errorf("the run xytspow3412oi must be in the 'active' state. Current state is deleted"),
			httpStatus:          http.StatusBadRequest,
			artifactURL:         "gs://bucketName/valid",
			deleteArtifact:      true,
			expectedRunRespJSON: `{}`,
		},
		{
			name:                "ID not exist",
			idRun:               "unknownId",
			expectedRespJSON:    RunDoesntExist,
			expectedError:       fmt.Errorf("run with id=unknownId not found"),
			httpStatus:          http.StatusNotFound,
			artifactURL:         "gs://bucketName/valid",
			deleteArtifact:      true,
			expectedRunRespJSON: `{}`,
		},
		{
			name:                "Artifact Deletion Failed",
			idRun:               "abcdefg1234",
			expectedRespJSON:    `{}`,
			expectedError:       fmt.Errorf("failed to Delete Artifact"),
			httpStatus:          http.StatusOK,
			artifactURL:         "gs://bucketName/invalid",
			deleteArtifact:      true,
			expectedRunRespJSON: `{}`,
		},
		{
			name:                "Delete without URL Valid",
			idRun:               "abcdefg1234",
			expectedRespJSON:    `{}`,
			expectedError:       nil,
			httpStatus:          http.StatusOK,
			artifactURL:         "",
			deleteArtifact:      true,
			expectedRunRespJSON: RunSuccessDeleteRun,
		},
		{
			name:                "Delete without URL Invalid",
			idRun:               "abcdefg1234",
			expectedRespJSON:    `{}`,
			expectedError:       fmt.Errorf("failed to Delete Artifact"),
			httpStatus:          http.StatusOK,
			artifactURL:         "",
			deleteArtifact:      true,
			expectedRunRespJSON: RunFailedDeleteRun,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/2.0/mlflow/runs/delete", func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(tc.httpStatus)
				_, err := w.Write([]byte(tc.expectedRespJSON))
				require.NoError(t, err)
			})
			mux.HandleFunc("/api/2.0/mlflow/runs/get", func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(tc.httpStatus)
				_, err := w.Write([]byte(tc.expectedRunRespJSON))
				require.NoError(t, err)
			})

			server := httptest.NewServer(mux)
			defer server.Close()

			artifactServiceMock := mocks.Service{}
			client := mlflowService{
				API:             server.Client(),
				ArtifactService: &artifactServiceMock,
				Config:          Config{TrackingURL: server.URL},
			}

			artifactServiceMock.
				On("DeleteArtifact", context.Background(), "gs://bucketName/invalid").
				Return(fmt.Errorf("failed to Delete Artifact"))
			artifactServiceMock.
				On("DeleteArtifact", context.Background(), "gs://bucketName/valid").
				Return(nil)
			errAPI := client.DeleteRun(context.Background(), tc.idRun, tc.artifactURL, tc.deleteArtifact)
			assert.Equal(t, tc.expectedError, errAPI)
		})
	}
}
