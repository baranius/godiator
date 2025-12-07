// Package tests provides test cases for the godiator library.
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

// Core Handler Test Suite
func TestRunHandlerCoreTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerCoreTestSuite))
}

func (s *HandlerCoreTestSuite) TestHandlerRegisteryActions() {
	core.AddHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	h, ok := core.GetHandler[samples.MyRequest, samples.MyResponse]()

	s.True(ok)
	s.NotNil(h)

	core.RemoveHandler[samples.MyRequest]()

	h, ok = core.GetHandler[samples.MyRequest, samples.MyResponse]()

	s.False(ok)
	s.Nil(h)
}

// Core Subscriber Test Suite
type SubscriberCoreTestSuite struct {
	suite.Suite
}

func TestRunSubscriberCoreTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriberCoreTestSuite))
}

func (s *SubscriberCoreTestSuite) TestSubscriberRegisteryActions() {
	core.AddSubscriber(&samples.MySubscriptionHandler[samples.MySubscriptionRequest]{})

	subscribers := core.GetSubscribers[samples.MySubscriptionRequest]()

	s.NotNil(subscribers)
	s.NotEmpty(subscribers)

	core.RemoveSubscriber[samples.MySubscriptionRequest]()

	subscribers = core.GetSubscribers[samples.MySubscriptionRequest]()

	s.Empty(subscribers)
}

// Core Pipeline Test Suite
type PipelineCoreTestSuite struct {
	suite.Suite
}

func TestRunPipelineCoreTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineCoreTestSuite))
}

func (s *PipelineCoreTestSuite) TestPipelineRegisteryActions() {
	pipeline := &samples.LoggingPipeline{}

	core.AddPipeline(pipeline)

	pipelines := core.GetPipelines()

	s.NotNil(pipelines)
	s.NotEmpty(pipelines)

	core.ClearPipelines()

	pipelines = core.GetPipelines()

	s.Empty(pipelines)
}
