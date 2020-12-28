package pipelines

import (
	"github.com/baranx/godiatr/examples/handler"
	"github.com/baranx/godiatr/godiatr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ValidationPipelineTestSuite struct {
	suite.Suite
	g godiatr.IGodiatr
}

func TestValidationPipelineSuite(t *testing.T){
	suite.Run(t, new(ValidationPipelineTestSuite))
}

func (s *ValidationPipelineTestSuite) SetupTest(){
	s.g = godiatr.GetInstance()

	s.g.RegisterPipeline(&ValidationPipeline{})

	s.g.Register(&handler.SampleRequest{}, handler.NewSampleHandler)
}

func (s *ValidationPipelineTestSuite) Test() {
	request := &handler.SampleRequest{}

	response, err := s.g.Send(request)

	assert.Nil(s.T(), response)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "PayloadString_should_not_be_null", err.Error())
}