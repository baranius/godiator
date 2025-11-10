// Package pipeline provides base implementations for the Pipeline interface.
// Pipelines are middleware components that can be used to implement cross-cutting
// concerns like logging, validation, or authentication in the mediator pattern.
package pipeline

import (
	"errors"

	"github.com/baranius/godiator/core/interfaces"
)

var _ interfaces.Pipeline = (*BasePipeline)(nil)

// BasePipeline is a default implementation of the Pipeline interface.
// It provides a mechanism to chain pipelines together and delegate
// request handling to the next pipeline in the chain.
type BasePipeline struct {
	nextPipeline interfaces.Pipeline
}

// Next returns the next pipeline in the chain.
//
// Returns:
//   - interfaces.Pipeline: The next pipeline, or nil if there is no next pipeline
func (p *BasePipeline) Next() interfaces.Pipeline {
	return p.nextPipeline
}

// SetNext sets the next pipeline in the chain.
//
// Parameters:
//   - nextPipeline: The next pipeline to set
func (p *BasePipeline) SetNext(nextPipeline interfaces.Pipeline) {
	p.nextPipeline = nextPipeline
}

// Handle processes the request and delegates to the next pipeline in the chain.
// This method should be overridden by custom pipelines to implement specific behavior.
//
// Parameters:
//   - request: The request object to process
//   - params: Optional additional parameters passed to the pipeline
//
// Returns:
//   - any: The response from the next pipeline or handler
//   - error: An error if processing fails
func (p *BasePipeline) Handle(request any, params ...any) (any, error) {
	return nil, errors.New("handle_method_not_implemented")
}
