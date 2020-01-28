package notify

import (
	"fmt"
	"github.com/baranx/godiatr/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MockingNotifyTestSuite struct {
	suite.Suite
	mockGodiatr mock.MockGodiatr
	controller MockNotificationController
}

func TestMockingNotifySuite(t *testing.T){
	suite.Run(t, new(MockingNotifyTestSuite))
}

func (s *MockingNotifyTestSuite) SetupTest(){
	s.mockGodiatr = mock.MockGodiatr{}
	s.controller = MockNotificationController{g: &s.mockGodiatr}
}

func (s *MockingNotifyTestSuite) Test_Godiatr_Notify() {
	// Given
	id := 1

	s.mockGodiatr.OnNotify = func(request interface{}, params ...interface{}){
		req := request.(*MockNotificationRequest)
		responseString := fmt.Sprintf("Requested id : %v", *req.Id)

		fmt.Print(responseString)
	}

	// When
	response, err := s.controller.GetItem(&id)

	// Then
	assert.Nil(s.T(), response)
	assert.Nil(s.T(), err)
}