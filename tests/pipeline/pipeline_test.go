package tests

import (
	"testing"

	"github.com/baranius/godiator/pipeline"
	"github.com/stretchr/testify/suite"
)

type PipelineTestSuite struct {
	suite.Suite
}

func TestRunPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}

func (s *PipelineTestSuite) TestPipelineActions() {
	pipeline := &pipeline.BasePipeline{}
	s.Suite.NotNil(pipeline)

	nextPipeline := &pipeline.BasePipeline{}
	pipeline.SetNext(nextPipeline)

	s.Suite.Equal(nextPipeline, pipeline.Next())

	result, err := pipeline.Handle(nil, nil...)
	s.Suite.Error(err)
	s.Suite.Nil(result)
}
