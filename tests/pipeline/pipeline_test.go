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
	firstPipeline := &pipeline.BasePipeline{}
	s.NotNil(firstPipeline)

	nextPipeline := &pipeline.BasePipeline{}
	firstPipeline.SetNext(nextPipeline)
	s.Equal(nextPipeline, firstPipeline.Next())

	result, err := firstPipeline.Handle(nil, nil...)
	s.Error(err)
	s.Nil(result)
}
