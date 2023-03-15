package mlflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gojek/mlp/api/pkg/gcs"
)

type MlflowService interface {
	SearchRunForExperiment(idExperiment string) (SearchRunsResponse, error)
	SearchRunData(idRun string) (SearchRunResponse, error)
	DeleteExperiment(idExperiment string) error
	DeleteRun(idRun string) error
}

type mlflowClient struct {
	Api        *http.Client
	GcsService gcs.GcsService
	Config     Config
}

func NewMlflowClient(httpClient *http.Client, config Config, gcsService gcs.GcsService) *mlflowClient {
	return &mlflowClient{
		Api:        httpClient,
		Config:     config,
		GcsService: gcsService,
	}
}

func (mfc *mlflowClient) httpCall(method string, url string, body []byte, response interface{}) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if method == "POST" {
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := mfc.Api.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		// Convert response body to Error Message struct
		var errMessage DeleteExperimentErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errMessage); err != nil {
			return err
		}
		return fmt.Errorf(errMessage.Message)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return err
		}
	}

	return nil
}

func (mfc *mlflowClient) SearchRunForExperiment(idExperiment string) (SearchRunsResponse, error) {
	// Search related runs for an experiment id
	var responseObject SearchRunsResponse

	searchRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/search", mfc.Config.TrackingURL)

	input := SearchRunRequest{ExperimentId: []string{idExperiment}}
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return responseObject, err
	}

	err = mfc.httpCall("POST", searchRunURL, jsonInput, &responseObject)
	if err != nil {
		return responseObject, err
	}

	return responseObject, nil
}

func (mfc *mlflowClient) SearchRunData(idRun string) (SearchRunResponse, error) {
	// Creating Output Format for Run Detail
	var runResponse SearchRunResponse
	getRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/get?run_id=%s", mfc.Config.TrackingURL, idRun)

	err := mfc.httpCall("GET", getRunURL, nil, &runResponse)
	if err != nil {
		return runResponse, err
	}
	return runResponse, nil
}

func (mfc *mlflowClient) DeleteExperiment(idExperiment string) error {

	relatedRunId, err := mfc.SearchRunForExperiment(idExperiment)
	if err != nil {
		return err
	}
	// Error Handling, keep runId that failed to delete
	// TODO: What should i do with this?
	var deletedRunId []string
	var failDeletedRunId []string
	for _, run := range relatedRunId.RunsData {
		err = mfc.DeleteRun(run.Info.RunId, false)
		if err != nil {
			failDeletedRunId = append(failDeletedRunId, run.Info.RunId)
			// return err
		} else {
			deletedRunId = append(deletedRunId, run.Info.RunId)
		}
	}

	if len(relatedRunId.RunsData) > 0 {
		// the [5:] is to remove the "gs://" on the artifact uri
		// ex : gs://bucketName/path → bucketName/path
		path := relatedRunId.RunsData[0].Info.ArtifactURI[5:]

		// This section is getting gcs path for a folder, from a run id
		// ex : bucketName/mlflow/idExperiment/runId/artefact  → bucketName/path/mlflow/idExperiment
		splitPath := strings.SplitN(path, "/", 4)
		folderPath := strings.Join(splitPath[0:3], "/")
		// deleting folder
		err = mfc.GcsService.DeleteArtifact(folderPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mfc *mlflowClient) DeleteRun(idRun string, deleteArtifact bool) error {
	// Creating Input Format for Delete run
	input := DeleteRunRequest{RunId: idRun}
	// HIT Delete Run API
	delRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/delete", mfc.Config.TrackingURL)

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return err
	}

	err = mfc.httpCall("POST", delRunURL, jsonInput, nil)
	if err != nil {
		return err
	}

	if deleteArtifact {
		runDetail, err := mfc.SearchRunData(idRun)
		if err != nil {
			return err
		}
		// the [5:] is to remove the "gs://" on the artifact uri
		// ex : gs://bucketName/path → bucketName/path
		err = mfc.GcsService.DeleteArtifact(runDetail.RunData.Info.ArtifactURI[5:])
		if err != nil {
			return err
		}

	}
	return nil
}
