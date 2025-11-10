package tests

import (
	"testing"

	"github.com/baranius/godiator"
	"github.com/baranius/godiator/samples"
	"github.com/stretchr/testify/suite"
)

type HandlerIntegrationTestSuite struct {
	suite.Suite
}

func TestHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerIntegrationTestSuite))
}

func (s *HandlerIntegrationTestSuite) TestHandlerExecutedSuccesfully() {
	// Given
	godiator.RegisterHandler(&samples.MyHandler[samples.MyRequest, samples.MyResponse]{})

	request := samples.MyRequest{
		Id: 1,
	}

	// When
	response, err := godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Nil(err)
	s.NotEmpty(response.Message)
	s.Equal("Processed successfully", response.Message)
}

func (s *HandlerIntegrationTestSuite) TestHandlerFailedExecution() {
	// Given
	godiator.RegisterHandler(&samples.MyFailedHandler[samples.MyFailedRequest, samples.MyFailedResponse]{})

	request := samples.MyFailedRequest{
		Reason: "Some failure reason",
	}

	// When
	response, err := godiator.Send[samples.MyFailedRequest, samples.MyFailedResponse](request, nil)

	// Then
	s.NotNil(err)
	s.Empty(response.Error)
}
