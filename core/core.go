package core

import (
	"reflect"

	"github.com/baranius/godiator/core/interfaces"
)

var (
	messageHandlers    = make(map[reflect.Type]interfaces.Handler[any, any])
	messageSubscribers = make(map[reflect.Type][]interfaces.Subscriber[any])
	messagePipelines   = make([]interfaces.Pipeline, 0)
)

// Wrapper for safe interfaces conversion
type handlerWrapper[TRequest any, TResponse any] struct {
	handler interfaces.Handler[TRequest, TResponse]
}

func (w *handlerWrapper[TRequest, TResponse]) Handle(request any, params ...any) (any, error) {
	return w.handler.Handle(request.(TRequest), params...)
}

// Wrapper for safe interfaces conversion
type subscriberWrapper[TRequest any] struct {
	subscriber interfaces.Subscriber[TRequest]
}

func (w *subscriberWrapper[TRequest]) Handle(request any, params ...any) {
	w.subscriber.Handle(request.(TRequest), params...)
}

// Adds the request model and handler pair
func AddHandler[TRequest any, TResponse any](handler interfaces.Handler[TRequest, TResponse]) {
	wrapper := &handlerWrapper[TRequest, TResponse]{handler}
	var request TRequest
	messageHandlers[reflect.TypeOf(request)] = wrapper
}

// Returns the handler for given request model
func GetHandler[TRequest any, TResponse any](request any) (*handlerWrapper[TRequest, TResponse], bool) {
	requestType := reflect.TypeOf(request)
	handler := messageHandlers[requestType]
	if handler == nil {
		return nil, false
	}

	return handler.(*handlerWrapper[TRequest, TResponse]), true
}

// Removes the handler for given request interfaces
func RemoveHandler[TRequest any]() {
	var request TRequest
	delete(messageHandlers, reflect.TypeOf(request))
}

// Registers the request model and handler(s) for fire & forget subscribers
func AddSubscriber[TRequest any](subscribers ...interfaces.Subscriber[TRequest]) {
	var request TRequest
	requestType := reflect.TypeOf(request)
	for _, s := range subscribers {
		wrapper := &subscriberWrapper[TRequest]{subscriber: s}
		if existing := messageSubscribers[requestType]; existing != nil {
			messageSubscribers[requestType] = append(existing, wrapper)
		} else {
			messageSubscribers[requestType] = []interfaces.Subscriber[any]{wrapper}
		}
	}
}

// Returns the subscriber(s) for given request model
func GetSubscribers[TRequest any](request TRequest) []subscriberWrapper[TRequest] {
	subscribers := messageSubscribers[reflect.TypeOf(request)]
	result := make([]subscriberWrapper[TRequest], 0)
	for _, sub := range subscribers {
		if wrapper, ok := sub.(*subscriberWrapper[TRequest]); ok {
			result = append(result, *wrapper)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// Removes the handler(s) for given request interfaces
func RemoveSubscriber[TRequest any](request TRequest) {
	delete(messageSubscribers, reflect.TypeOf(request))
}

// Adds Pipeline
func AddPipeline(p interfaces.Pipeline) {
	messagePipelines = append(messagePipelines, p)
}

// Returns all registered pipelines
func GetPipelines() []interfaces.Pipeline {
	return messagePipelines
}

// Removes all defined pipelines
func ClearPipelines() {
	messagePipelines = make([]interfaces.Pipeline, 0)
}
