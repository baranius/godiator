package samples

import (
	"testing"

	"github.com/baranius/godiator"
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
	godiator.RegisterHandler(&MyHandler[MyRequest, MyResponse]{})

	request := MyRequest{
		Id: 1,
	}

	// When
	response, err := godiator.Send[MyRequest, MyResponse](request, nil)

	// Then
	s.Suite.Nil(err)
	s.Suite.NotEmpty(response.Name)
	s.Suite.NotEmpty(response.Status)
}

func (s *HandlerIntegrationTestSuite) TestHandlerFailedExecution() {
	// Given
	godiator.RegisterHandler(&MyFailedHandler[MyRequest, MyResponse]{})

	request := MyRequest{
		Id: 1,
	}

	// When
	response, err := godiator.Send[MyRequest, MyResponse](request, nil)

	// Then
	s.Suite.NotNil(err)
	s.Suite.Empty(response)
}
