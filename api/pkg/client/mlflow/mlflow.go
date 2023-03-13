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
	SearchRunForExperiment(idExperiment string) (searchRunsResponse, error)
	SearchRunData(idRun string) (searchRunResponse, error)
	DeleteExperiment(idExperiment string) error
	DeleteRun(idRun string) error
}

type deleteExperimentRequest struct {
	ExperimentId string `json:"experiment_id" required:"true"`
}

type deleteRunRequest struct {
	RunId string `json:"run_id" required:"true"`
}

type deleteExperimentErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type searchRunRequest struct {
	ExperimentId []string `json:"experiment_ids" required:"true"`
}

type tagRun struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type dataRun struct {
	Tags []tagRun `json:"tags"`
}

type infoRun struct {
	RunId          string `json:"run_id"`
	ExperimentId   string `json:"experiment_id"`
	UserId         string `json:"user_id"`
	LifecycleStage string `json:"lifecycle_stage"`
	ArtifactURI    string `json:"artifact_uri"`
}
type runResponse struct {
	Info infoRun `json:"info"`
	Data dataRun `json:"data"`
}
type searchRunsResponse struct {
	RunsData []runResponse `json:"runs"`
}
type searchRunResponse struct {
	RunData runResponse `json:"run"`
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
		var errMessage deleteExperimentErrorResponse
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

func (mfc *mlflowClient) SearchRunForExperiment(idExperiment string) (searchRunsResponse, error) {
	// Search related runs for an experiment id
	var responseObject searchRunsResponse

	searchRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/search", mfc.Config.TrackingURL)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	input := searchRunRequest{ExperimentId: []string{idExperiment}}
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

func (mfc *mlflowClient) SearchRunData(idRun string) (searchRunResponse, error) {
	// Creating Output Format for Run Detail
	var runResponse searchRunResponse
	getRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/get?run_id=%s", mfc.Config.TrackingURL, idRun)

	err := mfc.httpCall("GET", getRunURL, nil, nil, &runResponse)
	if err != nil {
		return runResponse, err
	}
	return runResponse, nil
}

func (mfc *mlflowClient) DeleteExperiment(idExperiment string) error {
	// Creating Input Format for Delete experiment
	input := deleteExperimentRequest{ExperimentId: idExperiment}
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
	input := deleteRunRequest{RunId: idRun}
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
