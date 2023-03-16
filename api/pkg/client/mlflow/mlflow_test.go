package mlflow

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojek/mlp/api/pkg/gcs/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

var DeleteExperimentDoesntExist = `
{
    "error_code": "RESOURCE_DOES_NOT_EXIST",
    "message": "No Experiment with id=999 exists"
}`

var DeleteRunDoesntExist = `
{
    "error_code": "RESOURCE_DOES_NOT_EXIST",
    "message": "Run with id=unknownId not found"
}`

var DeleteRunAlreadyDeleted = `
{
    "error_code": "INVALID_PARAMETER_VALUE",
    "message": "The run xytspow3412oi must be in the 'active' state. Current state is deleted."
}`

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
							RunId:          "run-123",
							ExperimentId:   "1",
							UserId:         "root",
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
							RunId:          "run-456",
							ExperimentId:   "1",
							UserId:         "root",
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
			client := NewMlflowClient(server.Client(), Config{
				TrackingURL: server.URL,
			}, &mocks.GcsService{})

			resp, errAPI := client.SearchRunForExperiment(tc.idExperiment)
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
	}{
		{
			name:             "Valid Search",
			idRun:            "abcdefg1234",
			expectedRespJSON: RunSuccessJSON,
			expectedResponse: SearchRunResponse{
				RunData: RunResponse{
					Info: RunInfo{
						RunId:          "run-123",
						ExperimentId:   "1",
						UserId:         "root",
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
		},
		{
			name:             "No related runs",
			idRun:            "xytspow3412oi",
			expectedRespJSON: `{}`,
			expectedResponse: SearchRunResponse{},
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
			client := NewMlflowClient(server.Client(), Config{
				TrackingURL: server.URL,
			}, &mocks.GcsService{})

			resp, errAPI := client.SearchRunData(tc.idRun)

			assert.Equal(t, tc.expectedError, errAPI)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}

func TestMlflowClient_DeleteExperiment(t *testing.T) {
	tests := []struct {
		name             string
		idExperiment     string
		expectedRespJSON string
		expectedError    error
		httpStatus       int
	}{
		{
			name:             "Valid Experiment Deletion",
			idExperiment:     "1",
			expectedRespJSON: `{}`,
			expectedError:    nil,
			httpStatus:       http.StatusOK,
		},
		{
			name:             "ID not exist",
			idExperiment:     "999",
			expectedRespJSON: DeleteExperimentDoesntExist,
			expectedError:    fmt.Errorf("No Experiment with id=999 exists"),
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
			client := NewMlflowClient(server.Client(), Config{
				TrackingURL: server.URL,
			}, &mocks.GcsService{})

			errAPI := client.DeleteExperiment(tc.idExperiment)

			assert.Equal(t, tc.expectedError, errAPI)

		})
	}
}

func TestMlflowClient_DeleteRun(t *testing.T) {
	tests := []struct {
		name             string
		idRun            string
		expectedRespJSON string
		expectedError    error
		httpStatus       int
	}{
		{
			name:             "Valid Run Deletion",
			idRun:            "abcdefg1234",
			expectedRespJSON: `{}`,
			expectedError:    nil,
			httpStatus:       http.StatusOK,
		},
		{
			name:             "ID already deleted",
			idRun:            "xytspow3412oi",
			expectedRespJSON: DeleteRunAlreadyDeleted,
			expectedError:    fmt.Errorf("The run xytspow3412oi must be in the 'active' state. Current state is deleted."),
			httpStatus:       http.StatusBadRequest,
		},
		{
			name:             "ID not exist",
			idRun:            "unknownId",
			expectedRespJSON: DeleteRunDoesntExist,
			expectedError:    fmt.Errorf("Run with id=unknownId not found"),
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
			client := NewMlflowClient(server.Client(), Config{
				TrackingURL: server.URL,
			}, &mocks.GcsService{})

			errAPI := client.DeleteRun(tc.idRun, false)
			assert.Equal(t, tc.expectedError, errAPI)
		})
	}
}
