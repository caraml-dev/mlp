package mlflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type mlflowClient struct {
	Client *http.Client
	Config Config
}

type Config struct {
	TrackingURL string
}

type Mlflow interface {
	SearchRunForExperiment(idExperiment string) (SearchRunsResponse, error)
	SearchRunData(idRun string) (SearchRunResponse, error)
	DeleteExperiment(idExperiment string) error
	DeleteRun(idRun string) error
}

type DeleteExperimentRequest struct {
	ExperimentId string `json:"experiment_id" required:"true"`
}

type DeleteRunRequest struct {
	RunId string `json:"run_id" required:"true"`
}

type DeleteExperimentErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type SearchRunRequest struct {
	ExperimentId []string `json:"experiment_ids" required:"true"`
}

type TagRun struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type DataRun struct {
	Tags []TagRun `json:"tags"`
}

type InfoRun struct {
	RunId          string `json:"run_id"`
	ExperimentId   string `json:"experiment_id"`
	UserId         string `json:"user_id"`
	LifecycleStage string `json:"lifecycle_stage"`
	ArtifactURI    string `json:"artifact_uri"`
}
type RunResponse struct {
	Info InfoRun `json:"info"`
	Data DataRun `json:"data"`
}
type SearchRunsResponse struct {
	RunsData []RunResponse `json:"runs"`
}
type SearchRunResponse struct {
	RunData RunResponse `json:"run"`
}

func NewMlflowClient(httpClient *http.Client, config Config) *mlflowClient {
	return &mlflowClient{
		Client: httpClient,
		Config: config,
	}
}

func (mfc *mlflowClient) httpCall(method string, url string, headers map[string]string, body []byte, response interface{}) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := mfc.Client.Do(req)
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

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	input := SearchRunRequest{ExperimentId: []string{idExperiment}}
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return responseObject, err
	}

	err = mfc.httpCall("POST", searchRunURL, headers, jsonInput, &responseObject)
	if err != nil {
		return responseObject, err
	}

	return responseObject, nil
}

func (mfc *mlflowClient) SearchRunData(idRun string) (SearchRunResponse, error) {
	// Creating Output Format for Run Detail
	var runResponse SearchRunResponse
	getRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/get?run_id=%s", mfc.Config.TrackingURL, idRun)

	err := mfc.httpCall("GET", getRunURL, nil, nil, &runResponse)
	if err != nil {
		return runResponse, err
	}
	return runResponse, nil
}

func (mfc *mlflowClient) DeleteExperiment(idExperiment string) error {
	// Creating Input Format for Delete experiment
	input := DeleteExperimentRequest{ExperimentId: idExperiment}
	// HIT Delete Experiment API
	delExpURL := fmt.Sprintf("%s/api/2.0/mlflow/experiments/delete", mfc.Config.TrackingURL)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return err
	}

	err = mfc.httpCall("POST", delExpURL, headers, jsonInput, nil)
	if err != nil {
		return err
	}
	return nil
}

func (mfc *mlflowClient) DeleteRun(idRun string) error {
	// Creating Input Format for Delete run
	input := DeleteRunRequest{RunId: idRun}
	// HIT Delete Run API
	delRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/delete", mfc.Config.TrackingURL)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return err
	}

	err = mfc.httpCall("POST", delRunURL, headers, jsonInput, nil)
	if err != nil {
		return err
	}
	return nil
}
