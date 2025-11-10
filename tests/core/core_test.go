package tests

import (
	"testing"

	"github.com/baranius/godiator/core"
	"github.com/baranius/godiator/samples"
	"github.com/stretchr/testify/suite"
)

type HandlerCoreTestSuite struct {
	suite.Suite
}

func TestRunHandlerCoreTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerCoreTestSuite))
}

func (s *HandlerCoreTestSuite) TestHandlerActions() {
	handler := samples.MyHandler[samples.MyRequest, samples.MyResponse]{}

	core.AddHandler[samples.MyRequest, samples.MyResponse](&handler)

	h, ok := core.GetHandler[samples.MyRequest, samples.MyResponse]()

	s.True(ok)
	s.NotNil(h)

	core.RemoveHandler[samples.MyRequest]()

	h, ok = core.GetHandler[samples.MyRequest, samples.MyResponse]()

	s.False(ok)
	s.Nil(h)
}

type SubscriberCoreTestSuite struct {
	suite.Suite
}

func TestRunSubscriberCoreTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriberCoreTestSuite))
}

func (s *SubscriberCoreTestSuite) TestSubscriberActions() {
	subscriber := samples.MySubscriptionHandler[samples.MySubscriptionRequest]{}

	core.AddSubscriber[samples.MySubscriptionRequest](&subscriber)

	subscribers := core.GetSubscribers[samples.MySubscriptionRequest]()

	s.NotNil(subscribers)
	s.NotEmpty(subscribers)

	core.RemoveSubscriber[samples.MySubscriptionRequest]()

	subscribers = core.GetSubscribers[samples.MySubscriptionRequest]()

	s.Empty(subscribers)
}

type PipelineCoreTestSuite struct {
	suite.Suite
}

func TestRunPipelineCoreTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineCoreTestSuite))
}

func (s *PipelineCoreTestSuite) TestPipelineActions() {
	pipeline := &samples.LoggingPipeline{}

	core.AddPipeline(pipeline)

	pipelines := core.GetPipelines()

	s.NotNil(pipelines)
	s.NotEmpty(pipelines)

	core.ClearPipelines()

	pipelines = core.GetPipelines()

	s.Empty(pipelines)
}
