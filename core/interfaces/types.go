// Package interfaces defines the core interfaces used by the godiator library.
// These interfaces represent the contracts for handlers, subscribers, and pipelines
// that can be implemented by users to extend the functionality of the mediator.
package interfaces

// Handler represents a request/response handler in the mediator pattern.
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
type Handler[TRequest any, TResponse any] interface {
	Handle(request TRequest, params ...any) (TResponse, error)
}

// Subscriber represents a fire-and-forget handler in the mediator pattern.
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
type Subscriber[TRequest any] interface {
	Handle(request TRequest, params ...any)
}

// Pipeline represents a middleware component in the mediator pattern.
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
type Pipeline interface {
	Next() Pipeline
	SetNext(p Pipeline)
	Handle(request any, params ...any) (any, error)
}
