package samples

import (
	"testing"
	"time"

	"github.com/baranius/godiator"
	"github.com/baranius/godiator/register"
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
	req := MySubscriptionRequest{
		Id:     1,
		Status: true,
	}

	mySubscriber := &MySubscriptionHandler[MySubscriptionRequest]{IsHandlerExecuted: false}
	register.Subscriber(mySubscriber)

	// When
	godiator.Publish(req, nil)

	// Then
	time.Sleep(200 * time.Millisecond) // Wait for goroutine to complete
	s.Suite.True(mySubscriber.IsHandlerExecuted)
}

func (s *SubscriberIntegrationTestSuite) TestMultipleSubscriberExecution() {
	// Given
	req := MySubscriptionRequest{
		Id:     1,
		Status: true,
	}

	mySubscription := &MySubscriptionHandler[MySubscriptionRequest]{IsHandlerExecuted: false}
	myOtherSubscription := &MyOtherSubscriptionHandler[MySubscriptionRequest]{IsHandlerExecuted: false}

	register.Subscriber(mySubscription)
	register.Subscriber(myOtherSubscription)

	// When
	godiator.Publish(req, nil)

	// Then
	time.Sleep(200 * time.Millisecond) // Wait for goroutine to complete
	s.Suite.True(mySubscription.IsHandlerExecuted)
	s.Suite.True(myOtherSubscription.IsHandlerExecuted)
}
