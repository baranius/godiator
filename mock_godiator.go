package godiator

type OnSend func(request interface{}, params ...interface{}) (interface{}, error)
type OnPublish func(request interface{}, params ...interface{})

type MockGodiator struct {
	OnSend    OnSend
	OnPublish OnPublish
}

func (mock *MockGodiator) GetHandlerResponse(request interface{}) interface{} {
	return nil
}

func (mock *MockGodiator) RegisterPipeline(h IPipeline) {

}

func (mock *MockGodiator) Register(request interface{}, handler func() interface{}) {

}

func (mock *MockGodiator) RegisterSubscription(request interface{}, handler ...func() interface{}) {

}

func (mock *MockGodiator) Send(request interface{}, params ...interface{}) (interface{}, error) {
	return mock.OnSend(request, params...)
}

func (mock *MockGodiator) Publish(request interface{}, params ...interface{}) {
	mock.OnPublish(request, params...)
}
