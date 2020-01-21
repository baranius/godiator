package notification

import (
	"github.com/baranx/godiatr/godiatr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NotificationTestSuite struct {
	suite.Suite
	g godiatr.IGodiatr
}

func TestNotificationSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}

func (s *NotificationTestSuite) SetupTest() {
	s.g = godiatr.GetInstance()

	s.g.RegisterHandler(&NotificationCallerRequest{}, NewNotificationCallerHandler)

	s.g.RegisterNotificationHandler(&NotificationCallerRequest{}, NewNotification)
}

func (s *NotificationTestSuite) Test() {
	payload := "payload_string"

	request := &NotificationCallerRequest{PayloadString:&payload}

	response, err := s.g.Send(request)
	r := response.(*NotificationCallerResponse)

	assert.NotNil(s.T(), r)
	assert.Equal(s.T(), payload, *r.ResponseString)
	assert.Nil(s.T(), err)
}
