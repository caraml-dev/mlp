package pipelines

import "context"

type PipelineInterface interface {
	ListPipelines(ctx context.Context) (interface{}, error)
}
