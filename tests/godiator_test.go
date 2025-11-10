package tests

import (
	"testing"
	"time"

	"github.com/baranius/godiator"
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
	godiator.RegisterHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	// When
	response, err := godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Nil(err)
	s.Equal(samples.MyResponse{Message: "Processed successfully"}, response)
}

func (s *GodiatorTestSuite) TestGodiatorSend_WithPipeline() {
	// Given
	request := samples.MyRequest{Id: 2}
	pipeline := &samples.LoggingPipeline{}
	godiator.RegisterPipeline(pipeline)
	godiator.RegisterHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	// When
	response, err := godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Nil(err)
	s.Equal(samples.MyResponse{Message: "Processed successfully"}, response)
	s.Empty(pipeline.ErrorMessage)
	s.Equal(`request ({"Id":2}) | response ({"Message":"Processed successfully"})`, pipeline.LogMessage)
}

func (s *GodiatorTestSuite) TestGodiatorSend_WithMultiplePipeline() {
	// Given
	request := samples.MyRequest{Id: 3}
	firstPipeline := &samples.LoggingPipeline{}
	secondPipeline := &samples.LoggingPipeline{}
	godiator.RegisterPipeline(firstPipeline)
	godiator.RegisterPipeline(secondPipeline)
	godiator.RegisterHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	// When
	response, err := godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Nil(err)
	s.Equal(samples.MyResponse{Message: "Processed successfully"}, response)
	s.Empty(secondPipeline.ErrorMessage)
	s.Equal(`request ({"Id":3}) | response ({"Message":"Processed successfully"})`, secondPipeline.LogMessage)
	s.Empty(firstPipeline.ErrorMessage)
	s.Equal(`request ({"Id":3}) | response ({"Message":"Processed successfully"})`, firstPipeline.LogMessage)
}

type UnregisteredRequest struct {
	Value string
}

func (s *GodiatorTestSuite) TestGodiatorSend_HandlerNotFound() {
	// Given
	request := UnregisteredRequest{Value: "test"}

	// When
	response, err := godiator.Send[UnregisteredRequest, samples.MyResponse](request, nil)

	// Then
	s.NotNil(err)
	s.EqualError(err, `handler not found for "tests.UnregisteredRequest"`)
	s.Equal(samples.MyResponse{}, response)
}

func (s *GodiatorTestSuite) TestGodiatorPublish() {
	// Given
	request := samples.MySubscriptionRequest{Id: 1}
	subscriber := &samples.MySubscriptionHandler[samples.MySubscriptionRequest]{}
	godiator.RegisterSubscriber(subscriber)

	// When
	godiator.Publish[samples.MySubscriptionRequest](request, nil)

	// Then
	time.Sleep(200 * time.Millisecond)
	s.True(subscriber.IsHandlerExecuted)
}
