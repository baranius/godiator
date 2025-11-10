// Package samples provides example implementations of handlers and pipelines.
// These examples demonstrate how to use the godiator library in various scenarios.
package samples

import "errors"

// MyRequest represents a sample request structure.
type MyRequest struct {
	// Id is the identifier for the request.
	Id int
}

// MyResponse represents a sample response structure.
type MyResponse struct {
	// Message contains the response message.
	Message string
}

// MyHandler is a generic handler for processing requests and generating responses.
type MyHandler[TRequest MyRequest, TResponse MyResponse] struct{}

// Handle processes the given request and returns a response.
//
// Parameters:
//   - request: The request to process.
//   - params: Additional parameters for processing.
//
// Returns:
//   - MyResponse: The generated response.
//   - error: An error if processing fails.
func (mh *MyHandler[TRequest, TResponse]) Handle(request MyRequest, params ...any) (MyResponse, error) {
	response := MyResponse{
		Message: "Processed successfully",
	}
	return response, nil
}

// MyFailedRequest represents a sample failed request structure.
type MyFailedRequest struct {
	// Reason describes why the request failed.
	Reason string
}

// MyFailedResponse represents a sample failed response structure.
type MyFailedResponse struct {
	// Error contains the error message.
	Error string
}

// MyFailedHandler is a generic handler for processing failed requests.
type MyFailedHandler[TRequest MyFailedRequest, TResponse MyFailedResponse] struct{}

// Handle processes the given request and returns a response.
//
// Parameters:
//   - request: The request to process.
//   - params: Additional parameters for processing.
//
// Returns:
//   - MyFailedResponse: The generated response.
//   - error: An error if processing fails.
func (mh *MyFailedHandler[TRequest, TResponse]) Handle(request MyFailedRequest, params ...any) (MyFailedResponse, error) {
	return MyFailedResponse{}, errors.New("handler failed to process the request")
}
