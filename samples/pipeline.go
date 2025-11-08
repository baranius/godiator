package samples

import (
	"encoding/json"
	"fmt"

	"github.com/baranius/godiator/core"
)

type (
	LoggingPipeline struct {
		core.Pipeline
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
