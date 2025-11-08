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

	s.Suite.True(ok)
	s.Suite.NotNil(h)

	core.RemoveHandler[samples.MyRequest]()

	h, ok = core.GetHandler[samples.MyRequest, samples.MyResponse]()

	s.Suite.False(ok)
	s.Suite.Nil(h)
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

	s.Suite.NotNil(subscribers)
	s.Suite.NotEmpty(subscribers)

	core.RemoveSubscriber[samples.MySubscriptionRequest]()

	subscribers = core.GetSubscribers[samples.MySubscriptionRequest]()

	s.Suite.Empty(subscribers)
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

	s.Suite.NotNil(pipelines)
	s.Suite.NotEmpty(pipelines)

	core.ClearPipelines()

	pipelines = core.GetPipelines()

	s.Suite.Empty(pipelines)
}
