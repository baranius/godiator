package notify

import (
	"fmt"
	"github.com/baranx/godiatr/godiatr"
)

// GODIATR HANDLER
type (
	MockNotificationRequest struct {
		Id *int
	}

	MockNotificationResponse struct {
		ResponseString *string
	}

	MockNotificationHandler struct {
	}
)

func (h *MockNotificationHandler) Handle(request *MockNotificationRequest) {
	fmt.Print(request)
}

// CONTROLLER
type MockNotificationController struct {
	g godiatr.IGodiatr
}

func (c *MockNotificationController) GetItem(id *int) (interface{}, error) {
	request := &MockNotificationRequest{Id: id}

	c.g.Publish(request)

	return nil, nil
}