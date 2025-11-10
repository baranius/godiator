package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/baranius/godiator"
	"github.com/baranius/godiator/mockiator"
	"github.com/baranius/godiator/samples"
	"github.com/stretchr/testify/suite"
)

// Define an executer
func HandlerExecuter(id int) (samples.MyResponse, error) {
	request := samples.MyRequest{Id: id}

	return godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)
}

// Define Sample Subscriber godiator boilerplate
type (
	SubscriberRequest struct {
		Id int
	}
	SubscriberHandler[S1 SubscriberRequest] struct{}
)

func (h *SubscriberHandler[S1]) Handle(request SubscriberRequest, params ...any) {}

// Define an executer
func SubscriberExecuter(id int) {
	request := SubscriberRequest{Id: id}

	godiator.Publish[SubscriberRequest](request, nil)
}

// Execute Mocking Test
type MockiatorTestSuite struct {
	suite.Suite
}

func TestMockiatorTestSuite(t *testing.T) {
	suite.Run(t, new(MockiatorTestSuite))
}

func (s *MockiatorTestSuite) TestHandlerMocking() {
	// Given
	input := 10

	execution := mockiator.OnSend(func(request samples.MyRequest, params ...any) (samples.MyResponse, error) {
		fmt.Println(request.Id)
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
