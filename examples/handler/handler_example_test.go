package handler

import (
	"github.com/baranx/godiatr/godiatr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HandlerExampleTestSuite struct {
	suite.Suite
	g godiatr.IGodiatr
}

func TestHandlerExampleSuite(t *testing.T){
	suite.Run(t, new(HandlerExampleTestSuite))
}

func (s *HandlerExampleTestSuite) SetupTest(){
	s.g = godiatr.GetInstance()

	s.g.RegisterHandler(&SampleRequest{}, NewSampleHandler)
}

func (s *HandlerExampleTestSuite) Test() {
	value := "some_string"
	request := &SampleRequest{PayloadString: &value}

	result, err := s.g.Send(request)

	r := result.(*SampleResponse)
	assert.NotNil(s.T(), r)
	assert.Equal(s.T(), value, *r.ResultString)
	assert.Nil(s.T(), err)
}
