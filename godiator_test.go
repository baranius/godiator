package godiator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type GodiatorTestSuite struct {
	suite.Suite
}

func TestGodiatorSuite(t *testing.T) {
	suite.Run(t, new(GodiatorTestSuite))
}

func (s *GodiatorTestSuite) Test_GetInstance() {
	// When
	g := GetInstance()
	gdtr := g.(*godiator)

	// Then
	assert.IsType(s.T(), make(map[reflect.Type]func() interface{}), gdtr.handlers)
	assert.IsType(s.T(), make(map[reflect.Type][]func() interface{}), gdtr.notifications)
	assert.IsType(s.T(), make([]IPipeline, 0), gdtr.pipelines)
}

func (s *GodiatorTestSuite) Test_GetHandlerResponse() {
	// Given
	g := GetInstance()
	g.Register(&sampleRequest{}, newSampleHandler)

	// When
	responseObject := g.GetHandlerResponse(&sampleRequest{})

	// Then
	assert.NotNil(s.T(), responseObject)
}

func (s *GodiatorTestSuite) Test_Send_Should_Panic_When_Handler_Not_Found() {
	// Given
	g := GetInstance()

	g.Register(&sampleRequest{}, newSampleHandler)

	request := &failingRequest{}

	// When
	assert.Panics(s.T(), func() {
		g.Send(request)
	})
}

func (s *GodiatorTestSuite) Test_Send_Should_Panic_When_Handler_Not_Have_Handle_Func() {
	// Given
	g := GetInstance()

	g.Register(&failingRequest{}, newFailingHandler)

	request := &failingRequest{}

	// When
	assert.Panics(s.T(), func() {
		g.Send(request)
	})
}

func (s *GodiatorTestSuite) Test_Send_Should_Be_Executed() {
	// Given
	g := GetInstance()

	g.Register(&sampleRequest{}, newSampleHandler)

	sampleData := "test-string"
	request := &sampleRequest{PayloadString: &sampleData}

	// When
	response, err := g.Send(request)

	// Then
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), response)
	resp := response.(*sampleResponse)
	assert.Equal(s.T(), sampleData, *resp.ResultString)
}

func (s *GodiatorTestSuite) Test_Send_Should_Be_Executed_With_Pipeline() {
	// Given
	g := GetInstance()

	g.Register(&sampleRequest{}, newSampleHandler)
	g.RegisterPipeline(&validationPipeline{})

	request := &sampleRequest{}

	// When
	response, err := g.Send(request)

	// Then
	assert.Nil(s.T(), response)
	assert.NotNil(s.T(), err)
}

func (s *GodiatorTestSuite) Test_Subscriber_Should_Be_Executed() {
	// Given
	g := GetInstance()

	g.RegisterSubscription(&subscriberRequest{}, newSubscriberHandler)

	sampleData := "test-string"
	request := &subscriberRequest{PayloadString: &sampleData}

	// When
	g.Publish(request)
}
