// Test Suite for Mockiator
package tests

import (
	"testing"
	"time"

	"github.com/baranius/godiator"
	"github.com/baranius/godiator/mockiator"
	"github.com/baranius/godiator/samples"
	"github.com/stretchr/testify/suite"
)

// Define an executer to run the handler
func HandlerExecuter(id int) (samples.MyResponse, error) {
	request := samples.MyRequest{Id: id}

	return godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)
}

// Define Sample Subscriber godiator boilerplate
type (
	SubscriberRequest struct {
		Id int
	}
	SubscriberHandler[TRequest SubscriberRequest] struct{}
)

func (h *SubscriberHandler[TRequest]) Handle(request TRequest, params ...any) {}

// Define an executer to run the subscriber
func SubscriberExecuter(id int) {
	request := SubscriberRequest{Id: id}

	godiator.Publish(request, nil)
}

// Execute Mocking Test
type MockiatorTestSuite struct {
	suite.Suite
}

// Mockiator Test Suite
func TestMockiatorTestSuite(t *testing.T) {
	suite.Run(t, new(MockiatorTestSuite))
}

// Test Handler Mocking
func (s *MockiatorTestSuite) TestHandlerMocking() {
	// Given
	input := 10

	execution := mockiator.OnSend(func(request samples.MyRequest, params ...any) (samples.MyResponse, error) {
		return samples.MyResponse{
			Message: "Processed successfully",
		}, nil
	})

	// When
	resp, err := HandlerExecuter(input)

	// Then
	s.Nil(err)
	s.True(execution.IsCalled)
	s.Equal(1, execution.TimesCalled)
	s.Equal("Processed successfully", resp.Message)
}

// Test Subscriber Mocking
func (s *MockiatorTestSuite) TestSubscriberMocking() {
	// Given
	input := 10

	execution := mockiator.OnPublish(func(request SubscriberRequest, params ...any) {
		// Logic here
	})

	// When
	SubscriberExecuter(input)

	// Then
	time.Sleep(100 * time.Millisecond)
	s.True(execution.IsCalled)
	s.Equal(1, execution.TimesCalled)
}
