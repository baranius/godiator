package mock

import "github.com/baranx/godiatr/godiatr"

type OnSend func(request interface{}, params ...interface{}) (interface{}, error)
type OnNotify func(request interface{}, params ...interface{}) (interface{}, error)

type MockGodiatr struct {
	OnSend OnSend
	OnNotify OnNotify
}

func (mock *MockGodiatr) RegisterPipeline(h godiatr.IPipeline) {

}

func (mock *MockGodiatr) RegisterHandler(request interface{}, handler func()interface{}) {

}

func (mock *MockGodiatr) RegisterNotificationHandler(request interface{}, handler func()interface{}) {

}

func (mock *MockGodiatr) Send(request interface{}, params ...interface{}) (interface{}, error) {
	return mock.OnSend(request, params...)
}

func (mock *MockGodiatr) Notify(request interface{}, params ...interface{}) (interface{}, error) {
	return mock.OnNotify(request, params...)
}