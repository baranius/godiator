package pipelines

import (
	"errors"
	"github.com/baranx/godiatr/examples/handler"
	"github.com/baranx/godiatr/godiatr"
)

type ValidationPipeline struct {
	godiatr.Pipeline
}

func (p *ValidationPipeline) Handle(request interface{}, params ...interface{}) (interface{}, error) {
	r := request.(*handler.SampleRequest)

	if r.PayloadString == nil {
		return nil, errors.New("PayloadString_should_not_be_null")
	}

	return p.Next().Handle(request, params...)
}