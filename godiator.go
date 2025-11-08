package godiator

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/baranius/godiator/core"
	"github.com/baranius/godiator/core/interfaces"
)

func RegisterHandler[TRequest any, TResponse any](handler interfaces.Handler[TRequest, TResponse]) {
	core.AddHandler[TRequest, TResponse](handler)
}

func RegisterSubscriber[TRequest any](subscriber interfaces.Subscriber[TRequest]) {
	core.AddSubscriber[TRequest](subscriber)
}

func RegisterPipeline(pipeline interfaces.Pipeline) {
	core.AddPipeline(pipeline)
}

func UnregisterHandler[TRequest any]() {
	core.RemoveHandler[TRequest]()
}

func UnregisterSubscriber[TRequest any](subscriber interfaces.Subscriber[TRequest]) {
	core.RemoveSubscriber[TRequest]()
}

// Executes the related handler for given request type
func Send[TRequest any, TResponse any](request TRequest, params ...any) (TResponse, error) {
	handler, ok := core.GetHandler[TRequest, TResponse]()
	if !ok {
		var emptyResponse TResponse
		return emptyResponse, fmt.Errorf(`handler not found for "%s"`, reflect.TypeOf(request).String())
	}

	messagePipelines := core.GetPipelines()
	executionPipeline := &executionPipeline{
		wrapperFunc: handler.Handle,
	}

	var response any
	var err error

	if len(messagePipelines) > 0 {
		var firstPipeline interfaces.Pipeline
		for _, pipeline := range slices.Backward(messagePipelines) {
			if firstPipeline == nil {
				pipeline.SetNext(executionPipeline)
				firstPipeline = pipeline
			} else {
				pipeline.SetNext(firstPipeline)
				firstPipeline = pipeline
			}
		}
		response, err = firstPipeline.Handle(request, params...)
		return response.(TResponse), err
	} else {
		response, err := executionPipeline.Handle(request, params...)
		return response.(TResponse), err
	}
}

// Executes the related subscriber(s) for given request type
func Publish[TRequest any](request TRequest, params ...any) {
	subscribers := core.GetSubscribers[TRequest]()
	if len(subscribers) > 0 {
		for _, subscriber := range subscribers {
			go subscriber.Handle(request, params...)
		}
	} else {
		fmt.Printf(`handler not found for "%s" \n`, reflect.TypeOf(request).String())
	}
}
