package mlflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gojek/mlp/api/pkg/artifact"
)

type MlflowService interface {
	searchRunForExperiment(idExperiment string) (SearchRunsResponse, error)
	searchRunData(idRun string) (SearchRunResponse, error)
	DeleteExperiment(idExperiment string) error
	DeleteRun(idRun string) error
}

type mlflowService struct {
	Api             *http.Client
	ArtifactService artifact.ArtifactService
	Config          Config
}

func NewMlflowService(httpClient *http.Client, config Config, artifactService artifact.ArtifactService) *mlflowService {
	return &mlflowService{
		Api:             httpClient,
		Config:          config,
		ArtifactService: artifactService,
	}
}

func (mfs *mlflowService) httpCall(method string, url string, body []byte, response interface{}) error {
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

	resp, err := mfs.Api.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

func (mfs *mlflowService) searchRunsForExperiment(idExperiment string) (SearchRunsResponse, error) {
	// Search related runs for an experiment id
	var responseObject SearchRunsResponse

	searchRunsURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/search", mfs.Config.TrackingURL)

	input := SearchRunsRequest{ExperimentId: []string{idExperiment}}
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return responseObject, err
	}

	err = mfs.httpCall("POST", searchRunsURL, jsonInput, &responseObject)
	if err != nil {
		return responseObject, err
	}

	return responseObject, nil
}

func (mfs *mlflowService) searchRunData(idRun string) (SearchRunResponse, error) {
	// Creating Output Format for Run Detail
	var runResponse SearchRunResponse
	getRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/get?run_id=%s", mfs.Config.TrackingURL, idRun)

	err := mfs.httpCall("GET", getRunURL, nil, &runResponse)
	if err != nil {
		return runResponse, err
	}
	return runResponse, nil
}

func (mfs *mlflowService) DeleteExperiment(idExperiment string) error {

	relatedRunId, err := mfs.searchRunsForExperiment(idExperiment)
	if err != nil {
		return err
	}
	// Error Handling, when a runId failed to delete return error
	for _, run := range relatedRunId.RunsData {
		err = mfs.DeleteRun(run.Info.RunId)
		if err != nil {
			return fmt.Errorf("deletion failed for run_id %s for experiment id %s: %s", run.Info.RunId, idExperiment, err)
		}
	}

	return nil
}

func (mfs *mlflowService) DeleteRun(idRun string) error {
	// Creating Input Format for Delete run
	input := DeleteRunRequest{RunId: idRun}
	// HIT Delete Run API
	delRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/delete", mfs.Config.TrackingURL)

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return err
	}

	err = mfs.httpCall("POST", delRunURL, jsonInput, nil)
	if err != nil {
		return err
	}

	runDetail, err := mfs.searchRunData(idRun)
	if err != nil {
		return err
	}
	// the [5:] is to remove the "gs://" on the artifact uri
	// ex : gs://bucketName/path → bucketName/path
	err = mfs.ArtifactService.DeleteArtifact(runDetail.RunData.Info.ArtifactURI[5:])
	if err != nil {
		return err
	}
	return nil
}
