package tests

import (
	"testing"

	"github.com/baranius/godiator"
	"github.com/baranius/godiator/samples"
	"github.com/stretchr/testify/suite"
)

type PipelineIntegrationTestSuite struct {
	suite.Suite
}

func TestPipelineIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineIntegrationTestSuite))
}

func (s *PipelineIntegrationTestSuite) TestPipelineInterceptedSuccesfully() {
	// Given
	loggingPipeline := &samples.LoggingPipeline{}
	godiator.RegisterPipeline(loggingPipeline)

	request := samples.MyRequest{Id: 1}
	godiator.RegisterHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	// When
	response, err := godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Suite.NoError(err)
	s.Suite.NotNil(response)
	s.Suite.Equal(request.Id, response.Id)
	s.Suite.Equal("", loggingPipeline.ErrorMessage)
	s.Suite.Equal(`request ({"Id":1}) | response ({"Id":1,"Name":"John Doe","Status":"Unknown"})`, loggingPipeline.LogMessage)
}

func (s *PipelineIntegrationTestSuite) TestPipelineHandlesErrorSuccesfully() {
	// Given
	errorPipeline := &samples.LoggingPipeline{}
	godiator.RegisterPipeline(errorPipeline)

	request := samples.MyRequest{Id: 1}
	godiator.RegisterHandler(&samples.MyFailedHandler[samples.MyRequest, samples.MyResponse]{})

	// When
	response, err := godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Suite.Equal(0, response.Id)
	s.Suite.NotNil(err)
	s.Suite.Equal("handler failed to process the request", errorPipeline.ErrorMessage)
	s.Suite.Equal("", errorPipeline.LogMessage)
}
