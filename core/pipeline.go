package core

import (
	"errors"

	"github.com/baranius/godiator/core/interfaces"
)

// Ensure Pipeline implements interfaces.Pipeline interface
var _ interfaces.Pipeline = (*Pipeline)(nil)

// Base Pipeline for new pipeline definitions
type Pipeline struct {
	nextPipeline interfaces.Pipeline
}

// Returns the following pipeline
func (p *Pipeline) Next() interfaces.Pipeline {
	return p.nextPipeline
}

// Sets the following pipeline
func (p *Pipeline) SetNext(nextPipeline interfaces.Pipeline) {
	p.nextPipeline = nextPipeline
}

// Business handler. Should be overridden
func (p *Pipeline) Handle(request any, params ...any) (any, error) {
	return nil, errors.New("handle_method_not_implemented")
}
