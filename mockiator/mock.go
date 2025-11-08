package mockiator

import "github.com/baranius/godiator"

type mockHandler[TRequest any, TResponse any] struct {
	handlerFunc func(request TRequest, params ...any) (TResponse, error)
	IsCalled    bool
	TimesCalled int
}

func (m *mockHandler[TRequest, TResponse]) Handle(request TRequest, params ...any) (TResponse, error) {
	m.IsCalled = true
	m.TimesCalled++
	return m.handlerFunc(request, params)
}

func OnSend[TRequest any, TResponse any](handler func(request TRequest, params ...any) (TResponse, error)) *mockHandler[TRequest, TResponse] {
	h := mockHandler[TRequest, TResponse]{handlerFunc: handler}
	godiator.RegisterHandler[TRequest, TResponse](&h)
	return &h
}

type mockSubscriber[TRequest any] struct {
	handlerFunc func(request TRequest, params ...any)
	IsCalled    bool
	TimesCalled int
}

func (s *mockSubscriber[TRequest]) Handle(request TRequest, params ...any) {
	s.IsCalled = true
	s.TimesCalled++
	s.handlerFunc(request, params)
}

func OnPublish[TRequest any](handler func(request TRequest, params ...any)) *mockSubscriber[TRequest] {
	subs := mockSubscriber[TRequest]{handlerFunc: handler}
	godiator.RegisterSubscriber(&subs)
	return &subs
}
