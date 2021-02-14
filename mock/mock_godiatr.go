package mock

import "github.com/baranx/godiatr/godiatr"

type OnSend func(request interface{}, params ...interface{}) (interface{}, error)
type OnPublish func(request interface{}, params ...interface{})

type MockGodiatr struct {
	OnSend   OnSend
	OnPublish OnPublish
}

func (mock *MockGodiatr) GetHandler(request interface{}) interface{} {
	return nil
}

func (mock *MockGodiatr) GetHandlerResponse(request interface{}) interface{} {
	return nil
}

func (mock *MockGodiatr) RegisterPipeline(h godiatr.IPipeline) {

}

func (mock *MockGodiatr) Register(request interface{}, handler func()interface{}) {

}

func (mock *MockGodiatr) RegisterNotification(request interface{}, handler func()interface{}) {

}

func (mock *MockGodiatr) Send(request interface{}, params ...interface{}) (interface{}, error) {
	return mock.OnSend(request, params...)
}

func (mock *MockGodiatr) Publish(request interface{}, params ...interface{}){
	mock.OnPublish(request, params...)
}