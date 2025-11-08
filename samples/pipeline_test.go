package samples

import (
	"testing"

	"github.com/baranius/godiator"
	"github.com/baranius/godiator/register"
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
	register.Pipeline(loggingPipeline)

	request := MyRequest{Id: 1}
	register.Handler(&MyHandler[MyRequest, MyResponse]{})

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
	register.Pipeline(errorPipeline)

	request := MyRequest{Id: 1}
	register.Handler(&MyFailedHandler[MyRequest, MyResponse]{})

	// When
	response, err := godiator.Send[MyRequest, MyResponse](request, nil)

	// Then
	s.Suite.Equal(0, response.Id)
	s.Suite.NotNil(err)
	s.Suite.Equal("handler failed to process the request", errorPipeline.ErrorMessage)
	s.Suite.Equal("", errorPipeline.LogMessage)
}
