// Package core provides the internal registry and management for handlers, subscribers,
// and pipelines used by the godiator mediator implementation.
//
// This package maintains thread-unsafe global state and is intended to be used only
// through the public API in the godiator package. Direct use of this package is not
// recommended unless you need low-level control over the mediator's behavior.
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

// AddHandler registers a handler for a specific request and response type pair.
// Only one handler can be registered per request type. If a handler already exists
// for the request type, it will be replaced.
//
// Type parameters:
//   - TRequest: The request type that the handler will process
//   - TResponse: The response type that the handler will return
//
// Parameters:
//   - handler: The handler to register
func AddHandler[TRequest any, TResponse any](handler interfaces.Handler[TRequest, TResponse]) {
	wrapper := &handlerWrapper[TRequest, TResponse]{handler}
	var request TRequest
	messageHandlers[reflect.TypeOf(request)] = wrapper
}

// GetHandler returns a handler wrapper for the specified request and response types.
//
// Returns:
//   - *handlerWrapper[TRequest, TResponse]: The handler wrapper.
//   - bool: Indicates whether the handler was found.
func GetHandler[TRequest any, TResponse any]() (*handlerWrapper[TRequest, TResponse], bool) {
	var request TRequest
	requestType := reflect.TypeOf(request)
	handler := messageHandlers[requestType]
	if handler == nil {
		return nil, false
	}

	if handler, ok := handler.(*handlerWrapper[TRequest, TResponse]); ok {
		return handler, true
	}

	return nil, false
}

// RemoveHandler unregisters the handler for the specified request type.
//
// Type parameters:
//   - TRequest: The request type whose handler should be removed
func RemoveHandler[TRequest any]() {
	var request TRequest
	delete(messageHandlers, reflect.TypeOf(request))
}

// AddSubscriber registers one or more subscribers for a specific request type.
// Subscribers are executed asynchronously when Publish is called.
//
// Type parameters:
//   - TRequest: The request type that the subscribers will process
//
// Parameters:
//   - subscribers: The subscribers to register
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

// GetSubscribers returns a list of subscriber wrappers for the specified request type.
//
// Returns:
//   - []subscriberWrapper[TRequest]: The list of subscriber wrappers.
func GetSubscribers[TRequest any]() []subscriberWrapper[TRequest] {
	var request TRequest
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

// RemoveSubscriber unregisters all subscribers for the specified request type.
//
// Type parameters:
//   - TRequest: The request type whose subscribers should be removed
func RemoveSubscriber[TRequest any]() {
	var request TRequest
	delete(messageSubscribers, reflect.TypeOf(request))
}

// AddPipeline registers a pipeline that will be executed before handlers.
// Pipelines are executed in reverse order of registration (last registered, first executed).
// Use pipelines for cross-cutting concerns like logging, validation, or authentication.
//
// Parameters:
//   - p: The pipeline to register
func AddPipeline(p interfaces.Pipeline) {
	messagePipelines = append(messagePipelines, p)
}

// GetPipelines retrieves all registered pipelines.
//
// Returns:
//   - []interfaces.Pipeline: The list of registered pipelines
func GetPipelines() []interfaces.Pipeline {
	return messagePipelines
}

// ClearPipelines removes all registered pipelines.
func ClearPipelines() {
	messagePipelines = make([]interfaces.Pipeline, 0)
}
