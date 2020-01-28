package send

import (
	"fmt"
	"github.com/baranx/godiatr/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MockingSendTestSuite struct {
	suite.Suite
	mockGodiatr mock.MockGodiatr
	controller MockController
}

func TestMockingSendSuite(t *testing.T){
	suite.Run(t, new(MockingSendTestSuite))
}

func (s *MockingSendTestSuite) SetupTest(){
	s.mockGodiatr = mock.MockGodiatr{}
	s.controller = MockController{g: &s.mockGodiatr}
}

func (s *MockingSendTestSuite) Test_Godiatr_Send() {
	// Given
	id := 1
	request := &MockRequest{Id:&id}

	s.mockGodiatr.OnSend = func(request interface{}, params ...interface{}) (i interface{}, err error) {
		req := request.(*MockRequest)
		responseString := fmt.Sprintf("Requested id : %v", *req.Id)

		return &MockResponse{ResponseString: &responseString}, nil
	}

	// When
	response, err := s.controller.GetItem(&id)
	resp := response.(*MockResponse)

	// Then
	assert.NotNil(s.T(), response)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("Requested id : %v", *request.Id), *resp.ResponseString)
}