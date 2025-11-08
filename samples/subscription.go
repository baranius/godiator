package samples

type (
	MySubscriptionRequest struct {
		Id     int
		Status bool
	}

	MySubscriptionHandler[T1 MySubscriptionRequest] struct {
		IsHandlerExecuted bool
	}
	MyOtherSubscriptionHandler[T1 MySubscriptionRequest] struct {
		IsHandlerExecuted bool
	}
)

func (msh *MySubscriptionHandler[T1]) Handle(request MySubscriptionRequest, params ...any) {
	msh.IsHandlerExecuted = true
}

func (msh *MyOtherSubscriptionHandler[T1]) Handle(request MySubscriptionRequest, params ...any) {
	msh.IsHandlerExecuted = true
}
