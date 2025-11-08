package pipeline

import (
	"errors"

	"github.com/baranius/godiator/core/interfaces"
)

// Ensure Pipeline implements interfaces.Pipeline interface
var _ interfaces.Pipeline = (*BasePipeline)(nil)

// Base BasePipeline for new pipeline definitions
type BasePipeline struct {
	nextPipeline interfaces.Pipeline
}

// Returns the following pipeline
func (p *BasePipeline) Next() interfaces.Pipeline {
	return p.nextPipeline
}

// Sets the following pipeline
func (p *BasePipeline) SetNext(nextPipeline interfaces.Pipeline) {
	p.nextPipeline = nextPipeline
}

// Business handler. Should be overridden
func (p *BasePipeline) Handle(request any, params ...any) (any, error) {
	return nil, errors.New("handle_method_not_implemented")
}
