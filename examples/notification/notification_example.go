package notification

import (
	"fmt"
	"github.com/baranx/godiatr/godiatr"
)

type (
	NotificationCallerRequest struct {
		PayloadString *string
	}

	NotificationCallerResponse struct {
		ResponseString *string
	}

	NotificationCallerHandler struct {
		g godiatr.IGodiatr
	}
)

func NewNotificationCallerHandler() interface{} {
	return &NotificationCallerHandler{g:godiatr.GetInstance()}
}

func (h *NotificationCallerHandler) Handle(request *NotificationCallerRequest) (*NotificationCallerResponse, error) {
	if err := h.g.Notify(request); err != nil {
		return nil, err
	}

	return &NotificationCallerResponse{ResponseString: request.PayloadString}, nil
}


type Notification struct {

}

func NewNotification() interface{} {
	return &Notification{}
}

func (n *Notification) Handle(request interface{}, params ...interface{}) error {
	r := request.(*NotificationCallerRequest)
	fmt.Printf("Notification called with payload : '%v'", *r.PayloadString)

	return nil
}