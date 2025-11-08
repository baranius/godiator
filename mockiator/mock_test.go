package mockiator

import (
	"fmt"
	"testing"
	"time"

	"github.com/baranius/godiator"
	"github.com/stretchr/testify/suite"
)

// Define Sample godiator boilerplate
type (
	SampleRequest struct {
		Id int
	}

	SampleResponse struct {
		Id   int
		Name string
	}

	SampleHandler[S1 SampleRequest, S2 SampleResponse] struct{}
)

func (h *SampleHandler[S1, S2]) Handle(request SampleRequest, params ...any) (SampleResponse, error) {
	return SampleResponse{}, nil
}

// Define an executer
func HandlerExecuter(id int) (SampleResponse, error) {
	request := SampleRequest{Id: id}

	return godiator.Send[SampleRequest, SampleResponse](request, nil)
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

	execution := OnSend(func(request SampleRequest, params ...any) (SampleResponse, error) {
		fmt.Println(request.Id)
		return SampleResponse{
			Id:   10,
			Name: "John Doe",
		}, nil
	})

	// When
	resp, err := HandlerExecuter(input)

	// Then
	s.Suite.Nil(err)
	s.Suite.True(execution.IsCalled)
	s.Suite.Equal(1, execution.TimesCalled)
	s.Suite.Equal(input, resp.Id)
	s.Suite.Equal("John Doe", resp.Name)
}

func (s *MockiatorTestSuite) TestSubscriberMocking() {
	// Given
	input := 10

	execution := OnPublish(func(request SubscriberRequest, params ...any) {
		// Logic here
	})

	// When
	SubscriberExecuter(input)

	// Then
	time.Sleep(100 * time.Millisecond)
	s.Suite.True(execution.IsCalled)
	s.Suite.Equal(1, execution.TimesCalled)
}
