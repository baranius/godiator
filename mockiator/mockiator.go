// Package mockiator provides mock implementations for testing purposes.
// It includes mock handlers and subscribers for use in unit tests.
package mockiator

import "github.com/baranius/godiator"

// mockHandler is a mock implementation of the Handler interface.
// It is used for testing purposes to simulate handler behavior.
type mockHandler[TRequest any, TResponse any] struct {
	handlerFunc func(request TRequest, params ...any) (TResponse, error)
	IsCalled    bool
	TimesCalled int
}

// Handle processes the request using the mock handler function.
//
// Parameters:
//   - request: The request object to process.
//   - params: Additional parameters for processing.
//
// Returns:
//   - TResponse: The response from the handler function.
//   - error: An error if processing fails.
func (m *mockHandler[TRequest, TResponse]) Handle(request TRequest, params ...any) (TResponse, error) {
	m.IsCalled = true
	m.TimesCalled++
	return m.handlerFunc(request, params)
}

// OnSend creates and registers a mock handler for the specified request and response types.
//
// Parameters:
//   - handler: The handler function to use for processing requests.
//
// Returns:
//   - *mockHandler[TRequest, TResponse]: The created mock handler.
func OnSend[TRequest any, TResponse any](handler func(request TRequest, params ...any) (TResponse, error)) *mockHandler[TRequest, TResponse] {
	h := mockHandler[TRequest, TResponse]{handlerFunc: handler}
	godiator.RegisterHandler[TRequest, TResponse](&h)
	return &h
}

// mockSubscriber is a mock implementation of the Subscriber interface.
// It is used for testing purposes to simulate subscriber behavior.
type mockSubscriber[TRequest any] struct {
	handlerFunc func(request TRequest, params ...any)
	IsCalled    bool
	TimesCalled int
}

// Handle processes the request using the mock subscriber function.
//
// Parameters:
//   - request: The request object to process.
//   - params: Additional parameters for processing.
func (s *mockSubscriber[TRequest]) Handle(request TRequest, params ...any) {
	s.IsCalled = true
	s.TimesCalled++
	s.handlerFunc(request, params)
}

// OnPublish creates and registers a mock subscriber for the specified request type.
//
// Parameters:
//   - handler: The subscriber function to use for processing requests.
//
// Returns:
//   - *mockSubscriber[TRequest]: The created mock subscriber.
func OnPublish[TRequest any](handler func(request TRequest, params ...any)) *mockSubscriber[TRequest] {
	subs := mockSubscriber[TRequest]{handlerFunc: handler}
	godiator.RegisterSubscriber(&subs)
	return &subs
}
