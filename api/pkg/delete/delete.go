package delete

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
}
type runResponse struct {
	Info infoRun `json:"info"`
	Data dataRun `json:"data"`
}
type searchRunResponse struct {
	RunData []runResponse `json:"runs"`
}

// func httpCall()

func DeleteExperiment(idExperiment string, deleteArtefact bool) error {
	// Creating Input Format for Delete experiment
	input := deleteExperimentRequest{ExperimentId: idExperiment}
	// HIT Delete Experiment API
	// TODO: Create URL Based On Tracking URL
	trackingURL := "https://mlflow.d.ai.golabs.io"
	delExpURL := fmt.Sprintf("%s/api/2.0/mlflow/experiments/delete", trackingURL)

	jsonReq, err := json.Marshal(input)
	if err != nil {
		return err
	}
	resp, err := http.Post(delExpURL, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		// Convert response body to Error Message struct
		var errMessage deleteExperimentErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errMessage)
		if err != nil {
			// Handle the error
			fmt.Println("Error:", err)
			return err
		}
		return fmt.Errorf(errMessage.Message)
	}

	// Search For the available run
	relatedRunId, err := SearchRunForExperiment(idExperiment)
	if err != nil {
		return err
	}
	var deletedRunId []string
	var failDeletedRunId []string
	for _, runId := range relatedRunId {
		err = DeleteRun(runId)
		if err != nil {
			failDeletedRunId = append(failDeletedRunId, runId)
			// return err
		} else {
			deletedRunId = append(deletedRunId, runId)
		}
	}

	return nil
}

func SearchRunForExperiment(idExperiment string) ([]string, error) {
	// searchInput := []string{idExperiment}
	input := searchRunRequest{ExperimentId: []string{idExperiment}}
	var runID []string
	// HIT Delete Experiment API
	// TODO: Create URL Based On Tracking URL
	trackingURL := "https://mlflow.d.ai.golabs.io"

	searchRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/search", trackingURL)

	jsonReq, err := json.Marshal(input)
	if err != nil {
		return runID, err
	}

	runResp, err := http.Post(searchRunURL, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		return runID, err
	}

	defer runResp.Body.Close()
	fmt.Println(runResp.StatusCode)

	if !(runResp.StatusCode >= http.StatusOK && runResp.StatusCode < http.StatusMultipleChoices) {
		// Convert response body to Error Message struct
		var errMessage deleteExperimentErrorResponse
		// json.Unmarshal(bodyBytes, &errMessage)
		return runID, fmt.Errorf(errMessage.Message)
	}
	var responseObject searchRunResponse
	err = json.NewDecoder(runResp.Body).Decode(&responseObject)
	if err != nil {
		// Handle the error
		fmt.Println("Error:", err)
		return runID, err
	}

	for _, run := range responseObject.RunData {
		runID = append(runID, run.Info.RunId)
	}
	return runID, nil
}

func DeleteRun(idRun string) error {
	// Creating Input Format for Delete experiment
	input := deleteRunRequest{RunId: idRun}
	// HIT Delete Experiment API
	// TODO: Create URL Based On Tracking URL
	trackingURL := "https://mlflow.d.ai.golabs.io"
	delRunURL := fmt.Sprintf("%s/api/2.0/mlflow/runs/delete", trackingURL)

	jsonReq, err := json.Marshal(input)
	if err != nil {
		return err
	}
	resp, err := http.Post(delRunURL, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		// Convert response body to Error Message struct
		var errMessage deleteExperimentErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errMessage)
		if err != nil {
			// Handle the error
			fmt.Println("Error:", err)
			return err
		}
		return fmt.Errorf(errMessage.Message)
	}
	return nil
}

// func main() {
// 	data, err := SearchRunForExperiment("2220")
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("Experiment Deleted")
// 	}
// 	fmt.Println(data)
// 	fmt.Println(data[0])
// 	// fmt.Println(time.Now().Unix() * 1000)
// }
