// Package samples provides example implementations of handlers and pipelines.
// These examples demonstrate how to use the godiator library in various scenarios.
package samples

import (
	"encoding/json"
	"fmt"

	"github.com/baranius/godiator/pipeline"
)

type (
	// LoggingPipeline is a pipeline that logs the request and response.
	LoggingPipeline struct {
		// BasePipeline is the base pipeline that implements the pipeline interface.
		pipeline.BasePipeline
		// ErrorMessage is the error message if an error occurs.
		ErrorMessage string
		// LogMessage is the log message.
		LogMessage string
	}
)

// Handle implements the pipeline interface.
// Check the pipeline interface (https://github.com/baranius/godiator/blob/master/pipeline/pipeline.go) for more details.
func (p *LoggingPipeline) Handle(request any, params ...any) (any, error) {
	// Call the next pipeline in the chain.
	response, err := p.Next().Handle(request, params)

	// If an error occurs, return it.
	if err != nil {
		p.ErrorMessage = err.Error()
		return response, err
	}

	// Marshal the request and response to JSON.
	requestJson, _ := json.Marshal(request)
	responseJson, _ := json.Marshal(response)

	// Set the log message.
	p.LogMessage = fmt.Sprintf("request (%s) | response (%s)", string(requestJson), string(responseJson))

	// Return the response and error.
	return response, err
}
