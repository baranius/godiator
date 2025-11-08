package samples

import (
	"encoding/json"
	"fmt"

	"github.com/baranius/godiator/pipeline"
)

type (
	LoggingPipeline struct {
		pipeline.BasePipeline
		ErrorMessage string
		LogMessage   string
	}
)

func (p *LoggingPipeline) Handle(request any, params ...any) (any, error) {
	response, err := p.Next().Handle(request, params)

	if err != nil {
		p.ErrorMessage = err.Error()
		return response, err
	}

	requestJson, _ := json.Marshal(request)
	responseJson, _ := json.Marshal(response)

	p.LogMessage = fmt.Sprintf("request (%s) | response (%s)", string(requestJson), string(responseJson))

	return response, err
}
