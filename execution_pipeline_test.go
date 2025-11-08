package godiator

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExecutionPipelineTestSuite struct {
	suite.Suite
}

func TestExecutionPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionPipelineTestSuite))
}

func (s *ExecutionPipelineTestSuite) TestExecutionPipeline() {
	// Given
	exePipeline := &executionPipeline{
		wrapperFunc: func(request any, params ...any) (any, error) {
			return "wrapped response", nil
		},
	}

	// When
	response, err := exePipeline.Handle("test request", nil)

	// Then
	s.Suite.Nil(err)
	s.Suite.Equal("wrapped response", response)
}
