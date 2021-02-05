package processor_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list-bot/internal/processor"
)

type ProcessorSuite struct {
	suite.Suite
	processor *processor.Processor
}

func (s *ProcessorSuite) SetupTest() {
	s.processor = processor.New()
}

func (s *ProcessorSuite) TestEmptyMessage() {
	s.Require().Panics(func() { s.processor.Process("") })
}

func (s *ProcessorSuite) TestCreateSprint() {
	const timeDay = 24 * time.Hour
	begin := time.Unix(time.Now().Unix()-rand.Int63n(5000)+rand.Int63n(5000), 0).Truncate(timeDay)
	end := begin.Add(7 * timeDay)

	sprintHeader := fmt.Sprintf("%s - %s", s.timeToSprintDate(begin), s.timeToSprintDate(end))
	resp := s.processor.Process("/ns " + sprintHeader)
	s.Require().Equal(sprintHeader, resp)
}

func (s *ProcessorSuite) timeToSprintDate(d time.Time) string {
	return fmt.Sprintf("%02d.%02d", d.Day(), d.Month())
}

func TestProcessorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProcessorSuite))
}
