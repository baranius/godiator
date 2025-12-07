// Package samples provides example implementations of handlers and pipelines.
// These examples demonstrate how to use the godiator library in various scenarios.
package samples

type (
	// MySubscriptionRequest is the request type for the subscription.
	MySubscriptionRequest struct {
		Id     int
		Status bool
	}

	// MySubscriptionHandler is the handler type for the subscription.
	MySubscriptionHandler[TRequest MySubscriptionRequest] struct {
		IsHandlerExecuted bool
	}

	// MyOtherSubscriptionHandler is the handler type for the subscription.
	MyOtherSubscriptionHandler[TRequest MySubscriptionRequest] struct {
		IsHandlerExecuted bool
	}
)

// Handle implements the subscription interface.
func (msh *MySubscriptionHandler[TRequest]) Handle(request TRequest, params ...any) {
	msh.IsHandlerExecuted = true
}

// Handle implements the subscription interface.
func (msh *MyOtherSubscriptionHandler[TRequest]) Handle(request TRequest, params ...any) {
	msh.IsHandlerExecuted = true
}
