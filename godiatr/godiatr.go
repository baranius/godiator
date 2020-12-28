package godiatr

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	godiatr *Godiatr
	once    sync.Once
)

// Define Struct
type (
	IGodiatr interface {
		GetHandler(request interface{}) interface{}
		GetHandlerResponse(request interface{}) interface{}
		RegisterPipeline(h IPipeline)
		Register(request interface{}, handler func()interface{})
		RegisterNotification(request interface{}, handler func()interface{})
		Send(request interface{}, params ...interface{}) (interface{}, error)
		Notify(request interface{}, params ...interface{})
	}

	Godiatr struct {
		handlers      map[reflect.Type]func() interface{}
		notifications map[reflect.Type][]func() interface{}
		pipelines []IPipeline
	}
)

// Define Initialization Method
func GetInstance() *Godiatr {
	once.Do(func() {
		godiatr = &Godiatr{
			handlers: make(map[reflect.Type]func() interface{}),
			notifications: make(map[reflect.Type][]func() interface{}),
		}
	})
	return godiatr
}

func (m *Godiatr) GetHandler(request interface{}) interface{} {
	modelType := reflect.TypeOf(request)
	handlerFunc := m.handlers[modelType]
	if handlerFunc == nil {
		panic(fmt.Sprintf("Handler related to '%s' not found", modelType.Name()))
	}

	return handlerFunc()
}

func (m *Godiatr) GetHandlerResponse(request interface{}) interface{} {
	handler := m.GetHandler(request)
	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName("Handle")
	responseType := method.Type().Out(0)
	responseTypeKind := responseType.Kind()
	var pv interface{}

	if responseTypeKind == reflect.Slice {
		pv = reflect.MakeSlice(responseType, 0, 0).Interface()
	} else if responseTypeKind == reflect.Struct {
		pv = reflect.New(responseType).Interface()
	} else if responseTypeKind == reflect.Ptr {
		if responseType.Elem().Kind() == reflect.Struct {
			pv = reflect.New(responseType.Elem()).Interface()
		} else if responseType.Elem().Kind() == reflect.Slice {
			pv = reflect.MakeSlice(responseType.Elem(), 0, 0).Interface()
		}
	}
	return pv
}

// Apply Interface
func (m *Godiatr) Register(request interface{}, handler func() interface{}) {
	m.handlers[reflect.TypeOf(request)] = handler
}

func (m *Godiatr) RegisterPipeline(h IPipeline) {
	m.pipelines = append(m.pipelines, h)
}

func (m *Godiatr) RegisterNotification(request interface{}, handler func() interface{}) {
	handlers := m.notifications[reflect.TypeOf(request)]
	handlers = append(handlers, handler)
	m.notifications[reflect.TypeOf(request)] = handlers
}

func (m *Godiatr) Send(request interface{}, params ...interface{}) (interface{}, error) {
	// Initialize an anonymous handler
	runnerPipeline := new(runnerPipeline)
	runnerPipeline.mediator = m

	if len(m.pipelines) > 0 {
		// Loop through pipelines by reverse if exists and bind them to each other
		var mainPipeline IPipeline
		for i := len(m.pipelines) - 1; i >= 0; i-- {
			if i == len(m.pipelines)-1 {
				pipeline := m.pipelines[i]
				pipeline.SetNext(runnerPipeline)
				mainPipeline = m.pipelines[i]
			} else {
				m.pipelines[i].SetNext(mainPipeline)
				mainPipeline = m.pipelines[i]
			}
		}
		if mainPipeline != nil {
			// Call nested handlers w/+ given parameters
			return mainPipeline.Handle(request, params...)
		} else {
			panic("Pipeline error")
		}
	} else {
		// Call handler w/- pipeline if there is no pipeline
		return runnerPipeline.Handle(request, params...)
	}
}

func (m *Godiatr) Notify(request interface{}, params ...interface{}) {
	// Check if request is nil or not
	if request == nil {
		panic(fmt.Sprintf("Godiatr request should not be null!"))
	}

	// Retrieve handler by Request
	modelType := reflect.TypeOf(request)
	notificationFunctions := m.notifications[modelType]

	if notificationFunctions == nil {
		panic(fmt.Sprintf("Handler related to '%s' not found", modelType.Name()))
	}

	for _, notificationFunc := range notificationFunctions {
		handler := notificationFunc()

		// Initialize Handler
		handlerValue := reflect.ValueOf(handler)
		method := handlerValue.MethodByName("Handle")
		if method.Kind() != reflect.Func {
			panic(fmt.Sprintf("Handle function not found in %s", handlerValue.Type().Name()))
		}

		// Iterate parameters
		var inputs []reflect.Value
		inputs = append(inputs, reflect.ValueOf(request))

		for _, v := range params {
			inputs = append(inputs, reflect.ValueOf(v))
		}

		// Call required method with given parameters
		method.Call(inputs)
	}
}
