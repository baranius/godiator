package godiator

import (
	"errors"
)

// Failing Handler
type (
	failingRequest struct {
	}

	failingHandler struct {
	}
)

func newFailingHandler() interface{} {
	return &failingHandler{}
}

// Sample Handler
type (
	sampleRequest struct {
		PayloadString *string
	}

	sampleResponse struct {
		ResultString *string
	}

	sampleHandler struct {
	}
)

func newSampleHandler() interface{} {
	return &sampleHandler{}
}

func (h *sampleHandler) Handle(request *sampleRequest) (*sampleResponse, error) {
	return &sampleResponse{ResultString: request.PayloadString}, nil
}

// Sample Subscriber
type (
	subscriberRequest struct {
		PayloadString *string
	}

	subscriberHandler struct {
	}
)

func newSubscriberHandler() interface{} {
	return &subscriberHandler{}
}

func (h *subscriberHandler) Handle(request *subscriberRequest) {
	panic("dunno what to do")
}

// Sample Pipeline
type validationPipeline struct {
	Pipeline
}

func (p *validationPipeline) Handle(request interface{}, params ...interface{}) (interface{}, error) {
	r := request.(*sampleRequest)

	if r.PayloadString == nil {
		return nil, errors.New("PayloadString_should_not_be_null")
	}

	return p.Next().Handle(request, params...)
}
