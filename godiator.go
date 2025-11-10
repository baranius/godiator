// Package godiator provides a mediator pattern implementation for Go, enabling loosely coupled
// communication between components through request/response and publish/subscribe patterns.
//
// The mediator pattern helps reduce direct dependencies between components by introducing
// a central mediator that handles communication. This package supports:
//   - Request/Response pattern via handlers
//   - Publish/Subscribe pattern via subscribers
//   - Pipeline behavior for cross-cutting concerns
//   - Generic type safety for all operations
//
// Basic usage:
//
//	// Define a handler
//	type MyHandler struct{}
//	func (h *MyHandler) Handle(req MyRequest, params ...any) (MyResponse, error) {
//	    return MyResponse{Data: "result"}, nil
//	}
//
//	// Register and use
//	godiator.RegisterHandler[MyRequest, MyResponse](&MyHandler{})
//	response, err := godiator.Send[MyRequest, MyResponse](MyRequest{})
package godiator

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/baranius/godiator/core"
	"github.com/baranius/godiator/core/interfaces"
)

// RegisterHandler registers a handler for a specific request and response type pair.
// Only one handler can be registered per request type. If a handler already exists
// for the request type, it will be replaced.
//
// Type parameters:
//   - TRequest: The request type that the handler will process
//   - TResponse: The response type that the handler will return
//
// Example:
//
//	type GetUserHandler struct{}
//	func (h *GetUserHandler) Handle(req GetUserRequest, params ...any) (GetUserResponse, error) {
//	    return GetUserResponse{Name: "John"}, nil
//	}
//	godiator.RegisterHandler[GetUserRequest, GetUserResponse](&GetUserHandler{})
func RegisterHandler[TRequest any, TResponse any](handler interfaces.Handler[TRequest, TResponse]) {
	core.AddHandler[TRequest, TResponse](handler)
}

// RegisterSubscriber registers a subscriber for a specific request type.
// Multiple subscribers can be registered for the same request type.
// Subscribers are executed asynchronously when Publish is called.
//
// Type parameters:
//   - TRequest: The request type that the subscriber will process
//
// Example:
//
//	type EmailSubscriber struct{}
//	func (s *EmailSubscriber) Handle(req UserCreatedEvent, params ...any) {
//	    // Send email notification
//	}
//	godiator.RegisterSubscriber[UserCreatedEvent](&EmailSubscriber{})
func RegisterSubscriber[TRequest any](subscriber interfaces.Subscriber[TRequest]) {
	core.AddSubscriber[TRequest](subscriber)
}

// RegisterPipeline registers a pipeline that will be executed before handlers.
// Pipelines are executed in reverse order of registration (last registered, first executed).
// Use pipelines for cross-cutting concerns like logging, validation, or authentication.
//
// Example:
//
//	type LoggingPipeline struct {
//	    pipeline.BasePipeline
//	}
//	func (p *LoggingPipeline) Handle(request any, params ...any) (any, error) {
//	    log.Printf("Request: %v", request)
//	    return p.Next().Handle(request, params...)
//	}
//	godiator.RegisterPipeline(&LoggingPipeline{})
func RegisterPipeline(pipeline interfaces.Pipeline) {
	core.AddPipeline(pipeline)
}

// UnregisterHandler removes the registered handler for the specified request type.
// After unregistration, calls to Send with this request type will return an error.
//
// Type parameters:
//   - TRequest: The request type whose handler should be removed
//
// Example:
//
//	godiator.UnregisterHandler[GetUserRequest]()
func UnregisterHandler[TRequest any]() {
	core.RemoveHandler[TRequest]()
}

// UnregisterSubscriber removes all registered subscribers for the specified request type.
// After unregistration, calls to Publish with this request type will not execute any subscribers.
//
// Type parameters:
//   - TRequest: The request type whose subscribers should be removed
//
// Example:
//
//	godiator.UnregisterSubscriber[UserCreatedEvent]()
func UnregisterSubscriber[TRequest any]() {
	core.RemoveSubscriber[TRequest]()
}

// Send dispatches a request to its registered handler and returns the response.
// If pipelines are registered, they will be executed in reverse order of registration
// before the handler is invoked.
//
// Type parameters:
//   - TRequest: The request type to send
//   - TResponse: The expected response type
//
// Parameters:
//   - request: The request object to process
//   - params: Optional additional parameters passed through pipelines and to the handler
//
// Returns:
//   - TResponse: The response from the handler
//   - error: An error if the handler is not found or if processing fails
//
// Example:
//
//	response, err := godiator.Send[GetUserRequest, GetUserResponse](GetUserRequest{ID: 1})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(response.Name)
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

// Publish dispatches a request to all registered subscribers asynchronously.
// Each subscriber is executed in a separate goroutine, making this a fire-and-forget operation.
// If no subscribers are registered for the request type, a message is printed to stdout.
//
// Type parameters:
//   - TRequest: The request type to publish
//
// Parameters:
//   - request: The request object to publish to all subscribers
//   - params: Optional additional parameters passed to each subscriber
//
// Example:
//
//	godiator.Publish[UserCreatedEvent](UserCreatedEvent{UserID: 123, Email: "user@example.com"})
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
