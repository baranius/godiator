// Test Suite for Pipeline
package tests

import (
	"testing"

	"github.com/baranius/godiator/pipeline"
	"github.com/stretchr/testify/suite"
)

// Pipeline Test Suite
type PipelineTestSuite struct {
	suite.Suite
}

// Run Pipeline Test Suite
func TestRunPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}

// Test Pipeline Actions
func (s *PipelineTestSuite) TestPipelineActions() {
	firstPipeline := &pipeline.BasePipeline{}
	s.NotNil(firstPipeline)

	nextPipeline := &pipeline.BasePipeline{}
	firstPipeline.SetNext(nextPipeline)
	s.Equal(nextPipeline, firstPipeline.Next())

	result, err := firstPipeline.Handle(nil, nil...)
	s.Error(err)
	s.Nil(result)
}
