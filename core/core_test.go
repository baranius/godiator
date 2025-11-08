package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ---- HANDLER TESTS -----
type (
	SampleRequest struct {
		Id int
	}

	SampleResponse struct {
		Id   int
		Name string
	}

	SampleHandler[TRequest SampleRequest, TResponse SampleResponse] struct{}
)

func (sh *SampleHandler[TRequest, TResponse]) Handle(request SampleRequest, params ...any) (SampleResponse, error) {
	return SampleResponse{
		Id:   request.Id,
		Name: "Test",
	}, nil
}

type HandlerCoreTestSuite struct {
	suite.Suite
}

func TestRunHandlerCoreTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerCoreTestSuite))
}

func (s *HandlerCoreTestSuite) TestHandlerActions() {
	request := SampleRequest{Id: 1}
	handler := SampleHandler[SampleRequest, SampleResponse]{}

	AddHandler[SampleRequest, SampleResponse](&handler)

	h, ok := GetHandler[SampleRequest, SampleResponse](request)

	s.Suite.True(ok)
	s.Suite.NotNil(h)

	RemoveHandler[SampleRequest]()

	h, ok = GetHandler[SampleRequest, SampleResponse](request)

	s.Suite.False(ok)
	s.Suite.Nil(h)
}

// ---- SUBSCRIBER TESTS -----
type (
	SubscriberRequest struct {
		Id int
	}

	SubscriberHandler[T1 SubscriberRequest] struct{}
)

func (sh *SubscriberHandler[T1]) Handle(request SubscriberRequest, params ...any) {}

type SubscriberCoreTestSuite struct {
	suite.Suite
}

func TestRunSubscriberCoreTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriberCoreTestSuite))
}

func (s *SubscriberCoreTestSuite) TestSubscriberActions() {
	request := SubscriberRequest{Id: 1}
	subscriber := SubscriberHandler[SubscriberRequest]{}

	AddSubscriber[SubscriberRequest](&subscriber)

	subscribers := GetSubscribers(request)

	s.Suite.NotNil(subscribers)
	s.Suite.NotEmpty(subscribers)

	RemoveSubscriber(SubscriberRequest{})

	subscribers = GetSubscribers(request)

	s.Suite.Empty(subscribers)
}

// ---- PIPELINE TESTS -----
type (
	SamplePipeline struct {
		Pipeline
	}
)

func (p *SamplePipeline) Handle(request any, params ...any) (any, error) {
	return p.Next().Handle(request, params)
}

type PipelineCoreTestSuite struct {
	suite.Suite
}

func TestRunPipelineCoreTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineCoreTestSuite))
}

func (s *PipelineCoreTestSuite) TestPipelineActions() {
	pipeline := &SamplePipeline{}

	AddPipeline(pipeline)

	pipelines := GetPipelines()

	s.Suite.NotNil(pipelines)
	s.Suite.NotEmpty(pipelines)

	ClearPipelines()

	pipelines = GetPipelines()

	s.Suite.Empty(pipelines)
}
