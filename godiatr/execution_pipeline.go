package godiatr

import (
	"fmt"
	"reflect"
)

type executionPipeline struct {
	Pipeline
	gdtr *godiatr
}

func (p *executionPipeline) Handle(request interface{}, params ...interface{}) (interface{}, error) {
	// Check if request is nil or not
	if request == nil {
		panic(fmt.Sprintf("Godiatr request should not be null!"))
	}

	// Retrieve handler by Request
	handler := p.gdtr.GetHandler(request)

	// Initialize Handler
	handlerValue := reflect.ValueOf(handler)
	method := handlerValue.MethodByName("Handle")
	if method.Kind() != reflect.Func {
		panic(fmt.Sprintf("'Handle' function not found in %s", handlerValue.Type().Name()))
	}

	// Iterate parameters
	var inputs []reflect.Value
	inputs = append(inputs, reflect.ValueOf(request))

	for _, v := range params {
		inputs = append(inputs, reflect.ValueOf(v))
	}

	// Call required method with given parameters
	result := method.Call(inputs)

	// Return result
	if result[1].Interface() != nil {
		return nil, result[1].Interface().(error)
	}

	return result[0].Interface(), nil
}
