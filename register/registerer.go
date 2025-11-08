package register

import (
	"github.com/baranius/godiator/core"
	"github.com/baranius/godiator/core/interfaces"
)

// Registers the request model and handler pair for handlers
func Handler[TRequest any, TResponse any](handler interfaces.Handler[TRequest, TResponse]) {
	core.AddHandler[TRequest, TResponse](handler)
}

// Registers the request model and handler(s) for fire & forget subscribers
func Subscriber[TRequest any](subscribers ...interfaces.Subscriber[TRequest]) {
	core.AddSubscriber[TRequest](subscribers...)
}

// Registers the given interceptor as pipeline
func Pipeline(p interfaces.Pipeline) {
	core.AddPipeline(p)
}
