package send

import "github.com/baranx/godiatr/godiatr"

// GODIATR HANDLER
type (
	MockRequest struct {
		Id *int
	}

	MockResponse struct {
		ResponseString *string
	}

	MockHandler struct {
	}
)

func (h *MockHandler) Handle(request *MockRequest) (*MockResponse, error) {
	return nil, nil
}

// CONTROLLER
type MockController struct {
	g godiatr.IGodiatr
}

func (c *MockController) GetItem(id *int) (interface{}, error) {
	request := &MockRequest{Id: id}

	response, err := c.g.Send(request)

	return response, err
}
