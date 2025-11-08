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
	s.Suite.Nil(err)
	s.Suite.NotEmpty(response.Name)
	s.Suite.NotEmpty(response.Status)
}

func (s *HandlerIntegrationTestSuite) TestHandlerFailedExecution() {
	// Given
	godiator.RegisterHandler(&samples.MyFailedHandler[samples.MyRequest, samples.MyResponse]{})

	request := samples.MyRequest{
		Id: 1,
	}

	// When
	response, err := godiator.Send[samples.MyRequest, samples.MyResponse](request, nil)

	// Then
	s.Suite.NotNil(err)
	s.Suite.Empty(response)
}
