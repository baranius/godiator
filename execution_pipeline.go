package godiator

import (
	"github.com/baranius/godiator/core"
)

// Execution Pipeline is the last ring of the pipeline chain
type executionPipeline struct {
	core.Pipeline
	wrapperFunc func(request any, params ...any) (any, error)
}

func (ep *executionPipeline) Handle(request any, params ...any) (any, error) {
	return ep.wrapperFunc(request, params)
}
