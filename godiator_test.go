package godiator

import (
	"testing"
	"time"

	"github.com/baranius/godiator/samples"
	"github.com/stretchr/testify/suite"
)

type GodiatorTestSuite struct {
	suite.Suite
}

func TestGodiatorTestSuite(t *testing.T) {
	suite.Run(t, new(GodiatorTestSuite))
}

func (s *GodiatorTestSuite) TestGodiatorSend() {
	// Given
	request := samples.MyRequest{Id: 1}
	RegisterHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	// When
	response, err := Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Suite.Nil(err)
	s.Suite.Equal(samples.MyResponse{Id: 1, Name: "John Doe", Status: "Unknown"}, response)
}

func (s *GodiatorTestSuite) TestGodiatorSend_WithPipeline() {
	// Given
	request := samples.MyRequest{Id: 2}
	pipeline := &samples.LoggingPipeline{}
	RegisterPipeline(pipeline)
	RegisterHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	// When
	response, err := Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Suite.Nil(err)
	s.Suite.Equal(samples.MyResponse{Id: 2, Name: "John Doe", Status: "Unknown"}, response)
	s.Suite.Empty(pipeline.ErrorMessage)
	s.Suite.Equal(`request ({"Id":2}) | response ({"Id":2,"Name":"John Doe","Status":"Unknown"})`, pipeline.LogMessage)
}

func (s *GodiatorTestSuite) TestGodiatorSend_WithMultiplePipeline() {
	// Given
	request := samples.MyRequest{Id: 3}
	firstPipeline := &samples.LoggingPipeline{}
	secondPipeline := &samples.LoggingPipeline{}
	RegisterPipeline(firstPipeline)
	RegisterPipeline(secondPipeline)
	RegisterHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	// When
	response, err := Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Suite.Nil(err)
	s.Suite.Equal(samples.MyResponse{Id: 3, Name: "John Doe", Status: "Unknown"}, response)
	s.Suite.Empty(secondPipeline.ErrorMessage)
	s.Suite.Equal(`request ({"Id":3}) | response ({"Id":3,"Name":"John Doe","Status":"Unknown"})`, secondPipeline.LogMessage)
	s.Suite.Empty(firstPipeline.ErrorMessage)
	s.Suite.Equal(`request ({"Id":3}) | response ({"Id":3,"Name":"John Doe","Status":"Unknown"})`, firstPipeline.LogMessage)
}

type UnregisteredRequest struct {
	Value string
}

func (s *GodiatorTestSuite) TestGodiatorSend_HandlerNotFound() {
	// Given
	request := UnregisteredRequest{Value: "test"}

	// When
	response, err := Send[UnregisteredRequest, samples.MyResponse](request, nil)

	// Then
	s.Suite.NotNil(err)
	s.Suite.EqualError(err, `handler not found for "godiator.UnregisteredRequest"`)
	s.Suite.Equal(samples.MyResponse{}, response)
}

func (s *GodiatorTestSuite) TestGodiatorPublish() {
	// Given
	request := samples.MySubscriptionRequest{Id: 1}
	subscriber := &samples.MySubscriptionHandler[samples.MySubscriptionRequest]{}
	RegisterSubscriber(subscriber)

	// When
	Publish[samples.MySubscriptionRequest](request, nil)

	// Then
	time.Sleep(200 * time.Millisecond)
	s.Suite.True(subscriber.IsHandlerExecuted)
}
