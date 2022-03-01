package api

import "net/http"

type PipelinesController struct {
	*AppContext
}

func (c *PipelinesController) ListPipelines(r *http.Request, vars map[string]string, _ interface{}) *ApiResponse {
	ctx := r.Context()

	pipelines, err := c.PipelineService.ListPipelines(ctx)
	if err != nil {
		return InternalServerError(err.Error())
	}

	return Ok(pipelines)
}
