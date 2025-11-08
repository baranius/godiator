package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PipelineTestSuite struct {
	suite.Suite
}

func TestRunPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}

func (s *PipelineTestSuite) TestPipelineActions() {
	pipeline := &Pipeline{}
	s.Suite.NotNil(pipeline)

	nextPipeline := &Pipeline{}
	pipeline.SetNext(nextPipeline)

	s.Suite.Equal(nextPipeline, pipeline.Next())

	result, err := pipeline.Handle(nil, nil...)
	s.Suite.Error(err)
	s.Suite.Nil(result)
}
