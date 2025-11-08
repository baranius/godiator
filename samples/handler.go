package samples

import "errors"

type (
	MyRequest struct {
		Id int
	}

	MyResponse struct {
		Id     int
		Name   string
		Status string
	}

	MyHandler[TRequest MyRequest, TResponse MyResponse]       struct{}
	MyFailedHandler[TRequest MyRequest, TResponse MyResponse] struct{}
)

func (mh *MyHandler[TRequest, TResponse]) Handle(request MyRequest, params ...any) (MyResponse, error) {
	return MyResponse{
		Id:     request.Id,
		Name:   "John Doe",
		Status: "Unknown",
	}, nil
}

func (mh *MyFailedHandler[TRequest, TResponse]) Handle(request MyRequest, params ...any) (MyResponse, error) {
	return MyResponse{}, errors.New("handler failed to process the request")
}
