package mock

import "github.com/baranx/godiatr/godiatr"

type OnSend func(request interface{}, params ...interface{}) (interface{}, error)

type MockGodiatr struct {
	OnSend OnSend
}

func (mock *MockGodiatr) Register(request interface{}, handler func()interface{}) {

}

func (mock *MockGodiatr) RegisterPipeline(h godiatr.IPipeline) {

}

func (mock *MockGodiatr) Send(request interface{}, params ...interface{}) (interface{}, error) {
	return mock.OnSend(request, params...)
}
