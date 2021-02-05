package processor_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list-bot/internal/processor"
)

var errTestNewSprint = errors.New("test: failed to store sprint")

type ProcessorSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	mockRepository *MockRepository
	processor      *processor.Processor
}

func (s *ProcessorSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockRepository = NewMockRepository(s.mockCtrl)
	s.processor = processor.New(s.mockRepository)
}

func (s *ProcessorSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *ProcessorSuite) TestEmptyMessage() {
	s.Require().Panics(func() { s.processor.Process("") })
}

func (s *ProcessorSuite) TestCreateSprint() {
	const timeDay = 24 * time.Hour
	begin := time.Unix(time.Now().Unix()-rand.Int63n(5000)+rand.Int63n(5000), 0).UTC().Truncate(timeDay)
	end := begin.Add(7 * timeDay)

	s.mockRepository.EXPECT().CreateNewSprint(begin, end)

	sprintHeader := fmt.Sprintf("%s - %s", s.timeToSprintDate(begin), s.timeToSprintDate(end))
	resp := s.processor.Process("/ns " + sprintHeader)
	s.Require().Equal(sprintHeader, resp)
}

func (s *ProcessorSuite) TestCreateSprintInvalidFormat() {
	tests := []string{
		"01.02 -", "- 02.03", "a.04 - 05.06", "07.b - 08.09", "10.11 - c.12", "", "13.01 - 14.d",
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc, func() {
			resp := s.processor.Process("/ns " + tc)
			s.Require().Equal("Invalid format of new sprint. Should be `DD.MM - DD.MM` (e.g. `01.12 - 07.12`)", resp)
		})
	}
}

func (s *ProcessorSuite) TestCreateSprintFailedToStore() {
	s.mockRepository.EXPECT().CreateNewSprint(gomock.Any(), gomock.Any()).Return(errTestNewSprint)
	resp := s.processor.Process("/ns 01.02 - 03.04")
	s.Require().Equal("Oops! Failed to create new sprint. Try later.", resp)
}

func (s *ProcessorSuite) timeToSprintDate(d time.Time) string {
	return fmt.Sprintf("%02d.%02d", d.Day(), d.Month())
}

func TestProcessorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProcessorSuite))
}
