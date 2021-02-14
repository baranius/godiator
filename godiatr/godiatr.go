package godiatr

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	gdtr *godiatr
	once sync.Once
)

// Define Struct
type godiatr struct {
	handlers      map[reflect.Type]func() interface{}
	notifications map[reflect.Type][]func() interface{}
	pipelines     []IPipeline
}

// Define Initialization Method
func GetInstance() IGodiatr {
	once.Do(func() {
		gdtr = &godiatr{
			handlers:      make(map[reflect.Type]func() interface{}),
			notifications: make(map[reflect.Type][]func() interface{}),
		}
	})
	return gdtr
}

func (g *godiatr) GetHandler(request interface{}) interface{} {
	modelType := reflect.TypeOf(request)
	handlerFunc := g.handlers[modelType]
	if handlerFunc == nil {
		panic(fmt.Sprintf("Handler related to '%s' not found", modelType.Name()))
	}

	return handlerFunc()
}

func (g *godiatr) GetHandlerResponse(request interface{}) interface{} {
	handler := g.GetHandler(request)
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
func (g *godiatr) Register(request interface{}, handler func() interface{}) {
	g.handlers[reflect.TypeOf(request)] = handler
}

func (g *godiatr) RegisterPipeline(h IPipeline) {
	g.pipelines = append(g.pipelines, h)
}

func (g *godiatr) RegisterNotification(request interface{}, handler func() interface{}) {
	handlers := g.notifications[reflect.TypeOf(request)]
	handlers = append(handlers, handler)
	g.notifications[reflect.TypeOf(request)] = handlers
}

func (g *godiatr) Send(request interface{}, params ...interface{}) (interface{}, error) {
	// Initialize an anonymous handler
	runnerPipeline := new(executionPipeline)
	runnerPipeline.gdtr = g

	if len(g.pipelines) > 0 {
		// Loop through pipelines by reverse if exists and bind them to each other
		var mainPipeline IPipeline
		for i := len(g.pipelines) - 1; i >= 0; i-- {
			if i == len(g.pipelines)-1 {
				pipeline := g.pipelines[i]
				pipeline.SetNext(runnerPipeline)
				mainPipeline = g.pipelines[i]
			} else {
				g.pipelines[i].SetNext(mainPipeline)
				mainPipeline = g.pipelines[i]
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

func (g *godiatr) Publish(request interface{}, params ...interface{}) {
	// Check if request is nil or not
	if request == nil {
		panic(fmt.Sprintf("Godiatr request should not be null!"))
	}

	// Retrieve handler by Request
	modelType := reflect.TypeOf(request)
	notificationFunctions := g.notifications[modelType]

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
