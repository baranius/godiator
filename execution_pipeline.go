package godiator

import "github.com/baranius/godiator/pipeline"

// Execution Pipeline is the last ring of the pipeline chain
type executionPipeline struct {
	pipeline.BasePipeline
	wrapperFunc func(request any, params ...any) (any, error)
}

func (ep *executionPipeline) Handle(request any, params ...any) (any, error) {
	return ep.wrapperFunc(request, params)
}
