package godiator

import (
	"fmt"
	"reflect"
)

type executionPipeline struct {
	Pipeline
	g *godiator
}

func (p *executionPipeline) Handle(request interface{}, params ...interface{}) (interface{}, error) {
	// Check if request is nil or not
	if request == nil {
		panic(fmt.Sprintf("Godiator request should not be nil!"))
	}

	// Get handler by request
	handler := p.g.getHandler(request)

	// Get handle method
	method, err := g.getHandleMethod(handler)
	if err != nil {
		panic(err.Error())
	}

	// Iterate parameters
	var inputs []reflect.Value
	inputs = append(inputs, reflect.ValueOf(request))

	for _, v := range params {
		inputs = append(inputs, reflect.ValueOf(v))
	}

	// Call handle method with given parameters
	result := method.Call(inputs)

	// Return result
	if result[1].Interface() != nil {
		return nil, result[1].Interface().(error)
	}

	return result[0].Interface(), nil
}
