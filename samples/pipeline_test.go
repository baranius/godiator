package samples

import (
	"testing"

	"github.com/baranius/godiator"
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
	loggingPipeline := &LoggingPipeline{}
	godiator.RegisterPipeline(loggingPipeline)

	request := MyRequest{Id: 1}
	godiator.RegisterHandler(&MyHandler[MyRequest, MyResponse]{})

	// When
	response, err := godiator.Send[MyRequest, MyResponse](request, nil)

	// Then
	s.Suite.NoError(err)
	s.Suite.NotNil(response)
	s.Suite.Equal(request.Id, response.Id)
	s.Suite.Equal("", loggingPipeline.ErrorMessage)
	s.Suite.Equal(`request ({"Id":1}) | response ({"Id":1,"Name":"John Doe","Status":"Unknown"})`, loggingPipeline.LogMessage)
}

func (s *PipelineIntegrationTestSuite) TestPipelineHandlesErrorSuccesfully() {
	// Given
	errorPipeline := &LoggingPipeline{}
	godiator.RegisterPipeline(errorPipeline)

	request := MyRequest{Id: 1}
	godiator.RegisterHandler(&MyFailedHandler[MyRequest, MyResponse]{})

	// When
	response, err := godiator.Send[MyRequest, MyResponse](request, nil)

	// Then
	s.Suite.Equal(0, response.Id)
	s.Suite.NotNil(err)
	s.Suite.Equal("handler failed to process the request", errorPipeline.ErrorMessage)
	s.Suite.Equal("", errorPipeline.LogMessage)
}
