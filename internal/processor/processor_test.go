package processor_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list-bot/internal/processor"
)

type ProcessorSuite struct {
	suite.Suite
}

func (s *ProcessorSuite) TestEmptyMessage() {
	p := processor.New()
	s.Require().Panics(func() { p.Process("") })
}

func TestProcessorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProcessorSuite))
}
