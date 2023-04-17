package types

type DeleteExperimentRequest struct {
	ExperimentID string `json:"experiment_id" required:"true"`
}

type DeleteRunRequest struct {
	RunID string `json:"run_id" required:"true"`
}

type DeleteExperimentErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type SearchRunsRequest struct {
	ExperimentID []string `json:"experiment_ids" required:"true"`
}

type RunTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type RunData struct {
	Tags []RunTag `json:"tags"`
}

type RunInfo struct {
	RunID          string `json:"run_id"`
	ExperimentID   string `json:"experiment_id"`
	UserID         string `json:"user_id"`
	LifecycleStage string `json:"lifecycle_stage"`
	ArtifactURI    string `json:"artifact_uri"`
}
type RunResponse struct {
	Info RunInfo `json:"info"`
	Data RunData `json:"data"`
}
type SearchRunsResponse struct {
	RunsData []RunResponse `json:"runs"`
}
type SearchRunResponse struct {
	RunData RunResponse `json:"run"`
}

type Config struct {
	TrackingURL         string
	ArtifactServiceType string
}
