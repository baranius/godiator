package tests

import (
	"testing"
	"time"

	"github.com/baranius/godiator"
	"github.com/baranius/godiator/samples"
	"github.com/stretchr/testify/suite"
)

type SubscriberIntegrationTestSuite struct {
	suite.Suite
}

func TestSubscriberIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriberIntegrationTestSuite))
}

func (s *SubscriberIntegrationTestSuite) TestSubscriberExecution() {
	// Given
	req := samples.MySubscriptionRequest{
		Id:     1,
		Status: true,
	}

	mySubscriber := &samples.MySubscriptionHandler[samples.MySubscriptionRequest]{IsHandlerExecuted: false}
	godiator.RegisterSubscriber(mySubscriber)

	// When
	godiator.Publish(req, nil)

	// Then
	time.Sleep(200 * time.Millisecond) // Wait for goroutine to complete
	s.True(mySubscriber.IsHandlerExecuted)
}

func (s *SubscriberIntegrationTestSuite) TestMultipleSubscriberExecution() {
	// Given
	req := samples.MySubscriptionRequest{
		Id:     1,
		Status: true,
	}

	mySubscription := &samples.MySubscriptionHandler[samples.MySubscriptionRequest]{IsHandlerExecuted: false}
	myOtherSubscription := &samples.MyOtherSubscriptionHandler[samples.MySubscriptionRequest]{IsHandlerExecuted: false}

	godiator.RegisterSubscriber(mySubscription)
	godiator.RegisterSubscriber(myOtherSubscription)

	// When
	godiator.Publish(req, nil)

	// Then
	time.Sleep(200 * time.Millisecond) // Wait for goroutine to complete
	s.True(mySubscription.IsHandlerExecuted)
	s.True(myOtherSubscription.IsHandlerExecuted)
}
