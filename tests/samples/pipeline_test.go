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
	s.NoError(err)
	s.NotNil(response)
	s.Equal("", loggingPipeline.ErrorMessage)
	s.Equal(`request ({"Id":1}) | response ({"Message":"Processed successfully"})`, loggingPipeline.LogMessage)
}

func (s *PipelineIntegrationTestSuite) TestPipelineHandlesErrorSuccesfully() {
	// Given
	errorPipeline := &samples.LoggingPipeline{}
	godiator.RegisterPipeline(errorPipeline)

	request := samples.MyFailedRequest{Reason: "Some failure reason"}
	godiator.RegisterHandler(&samples.MyFailedHandler[samples.MyFailedRequest, samples.MyFailedResponse]{})

	// When
	response, err := godiator.Send[samples.MyFailedRequest, samples.MyFailedResponse](request, nil)

	// Then
	s.NotNil(err)
	s.Empty(response)
	s.Equal("handler failed to process the request", errorPipeline.ErrorMessage)
	s.Equal("", errorPipeline.LogMessage)
}
