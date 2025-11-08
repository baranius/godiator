package register

import (
	"testing"

	"github.com/baranius/godiator/core"
	"github.com/stretchr/testify/suite"
)

type RegistererTestSuite struct {
	suite.Suite
}

type MyRequest struct {
	Id int
}

type MyResponse struct {
	Id     int
	Name   string
	Status string
}

type MyHandler[TRequest MyRequest, TResponse MyResponse] struct{}

func (h *MyHandler[TRequest, TResponse]) Handle(request MyRequest, params ...any) (MyResponse, error) {
	return MyResponse{
		Id:     request.Id,
		Name:   "John Doe",
		Status: "Active",
	}, nil
}

type MySubscriptionRequest struct {
	Id int
}

type MySubscriptionHandler[TRequest MySubscriptionRequest] struct{}

func (h *MySubscriptionHandler[TRequest]) Handle(request MySubscriptionRequest, params ...any) {
	// Subscription handling logic
}

type MyPipeline struct {
	core.Pipeline
}

func (p *MyPipeline) Handle(request any, params ...any) (any, error) {
	return nil, nil
}

func TestRegistererTestSuite(t *testing.T) {
	suite.Run(t, new(RegistererTestSuite))
}

func (s *RegistererTestSuite) TestRegisterHandler() {
	// Given
	request := MyRequest{Id: 1}
	handler := &MyHandler[MyRequest, MyResponse]{}
	Handler(handler)

	// When
	registeredHandler, isExists := core.GetHandler[MyRequest, MyResponse](request)

	// Then
	s.True(isExists)
	s.NotNil(registeredHandler)
}

func (s *RegistererTestSuite) TestRegisterSubscription() {
	// Given
	request := MySubscriptionRequest{Id: 1}
	subscriber := &MySubscriptionHandler[MySubscriptionRequest]{}
	Subscriber(subscriber)

	// When
	registeredSubscribers := core.GetSubscribers(request)

	// Then
	s.NotNil(registeredSubscribers)
	s.Len(registeredSubscribers, 1)
}

func (s *RegistererTestSuite) TestRegisterPipeline() {
	// Given
	pipeline := &MyPipeline{}
	Pipeline(pipeline)

	// When
	registeredPipelines := core.GetPipelines()

	// Then
	s.NotNil(registeredPipelines)
	s.Len(registeredPipelines, 1)
}
