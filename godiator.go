package godiator

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	g    *godiator
	once sync.Once
)

// Define Struct
type godiator struct {
	handlers      map[reflect.Type]func() interface{}
	notifications map[reflect.Type][]func() interface{}
	pipelines     []IPipeline
}

// Init Singleton
func GetInstance() IGodiator {
	once.Do(func() {
		g = &godiator{
			handlers:      make(map[reflect.Type]func() interface{}),
			notifications: make(map[reflect.Type][]func() interface{}),
		}
	})
	return g
}

func (g *godiator) getHandleMethod(handler interface{}) (reflect.Value, error) {
	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName("Handle")
	if method.Kind() != reflect.Func {
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("'Handle' function not found in %s", handlerValue.Type().String()))
	}
	return method, nil
}

func (g *godiator) getHandler(request interface{}) interface{} {
	modelType := reflect.TypeOf(request)
	handlerFunc := g.handlers[modelType]
	if handlerFunc == nil {
		panic(fmt.Sprintf("Handler related to '%s' not found", modelType.String()))
	}

	return handlerFunc()
}

// Apply Interface
func (g *godiator) GetHandlerResponse(request interface{}) interface{} {
	handler := g.getHandler(request)
	method, err := g.getHandleMethod(handler)
	if err != nil {
		panic(err.Error())
	}
	responseType := method.Type().Out(0)
	responseTypeKind := responseType.Kind()
	var handlerResponse interface{}

	if responseTypeKind == reflect.Slice {
		handlerResponse = reflect.MakeSlice(responseType, 0, 0).Interface()
	} else if responseTypeKind == reflect.Struct {
		handlerResponse = reflect.New(responseType).Interface()
	} else if responseTypeKind == reflect.Ptr {
		if responseType.Elem().Kind() == reflect.Struct {
			handlerResponse = reflect.New(responseType.Elem()).Interface()
		} else if responseType.Elem().Kind() == reflect.Slice {
			handlerResponse = reflect.MakeSlice(responseType.Elem(), 0, 0).Interface()
		}
	}
	return handlerResponse
}

func (g *godiator) Register(request interface{}, handler func() interface{}) {
	g.handlers[reflect.TypeOf(request)] = handler
}

func (g *godiator) RegisterPipeline(h IPipeline) {
	g.pipelines = append(g.pipelines, h)
}

func (g *godiator) RegisterSubscription(request interface{}, handlers ...func() interface{}) {
	g.notifications[reflect.TypeOf(request)] = handlers
}

func (g *godiator) Send(request interface{}, params ...interface{}) (interface{}, error) {
	// Initialize an anonymous pipeline
	executionPipeline := new(executionPipeline)
	executionPipeline.g = g

	if len(g.pipelines) > 0 {
		// Reverse loop through pipelines if any. Bind them to each other
		var mainPipeline IPipeline
		for i := len(g.pipelines) - 1; i >= 0; i-- {
			if i == len(g.pipelines)-1 {
				pipeline := g.pipelines[i]
				pipeline.SetNext(executionPipeline)
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
		return executionPipeline.Handle(request, params...)
	}
}

func (g *godiator) publishWithRecover(handlerName string, method reflect.Value, inputs []reflect.Value) {
	defer func() {
		if r := recover(); r != nil {
			message := fmt.Sprintf("Notification Failed (%v) -> ", handlerName)
			fmt.Println(message, r)
		}
	}()

	method.Call(inputs)
}

func (g *godiator) Publish(request interface{}, params ...interface{}) {
	// Check if request is nil or not
	if request == nil {
		panic(fmt.Sprintf("Godiator request should not be nil!"))
	}

	// Retrieve handler by Request
	modelType := reflect.TypeOf(request)
	notificationFunctions := g.notifications[modelType]

	if notificationFunctions == nil {
		panic(fmt.Sprintf("Handler related to '%s' not found", modelType.String()))
	}

	for _, notificationFunc := range notificationFunctions {
		handler := notificationFunc()

		// Initialize Handler
		handlerValue := reflect.ValueOf(handler)
		method := handlerValue.MethodByName("Handle")
		if method.Kind() != reflect.Func {
			panic(fmt.Sprintf("Handle function not found in %s", handlerValue.Type().String()))
		}

		// Iterate parameters
		var inputs []reflect.Value
		inputs = append(inputs, reflect.ValueOf(request))

		for _, v := range params {
			inputs = append(inputs, reflect.ValueOf(v))
		}

		// Call with given params
		g.publishWithRecover(reflect.TypeOf(handler).Elem().String(), method, inputs)
	}
}
