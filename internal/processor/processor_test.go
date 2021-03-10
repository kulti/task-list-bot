package processor_test

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list-bot/internal/processor"
)

var (
	errTestNewSprint   = errors.New("test: failed to store sprint")
	errTestNewTask     = errors.New("test: failed to create task")
	errTestDoneTask    = errors.New("test: failed to done task")
	errTestBurnPoints  = errors.New("test: failed to burn task points")
	errTestCurrentList = errors.New("test: failed to list tasks")
)

type ProcessorSuite struct {
	suite.Suite
	mockCtrl  *gomock.Controller
	mockStore *MockStore
	processor *processor.Processor
}

func (s *ProcessorSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockStore = NewMockStore(s.mockCtrl)
	s.processor = processor.New(s.mockStore)
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

	s.mockStore.EXPECT().CreateNewSprint(begin, end)
	dump := faker.Sentence()
	s.mockStore.EXPECT().CurrentSprintDump().Return(dump, nil)

	sprintHeader := fmt.Sprintf("%s - %s", s.timeToSprintDate(begin), s.timeToSprintDate(end))
	resp := s.processor.Process("/ns " + sprintHeader)
	s.Require().Equal(dump, resp)
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
	s.mockStore.EXPECT().CreateNewSprint(gomock.Any(), gomock.Any()).Return(errTestNewSprint)
	resp := s.processor.Process("/ns 01.02 - 03.04")
	s.Require().Equal("Oops! Failed to create new sprint. Try later.", resp)
}

func (s *ProcessorSuite) timeToSprintDate(d time.Time) string {
	return fmt.Sprintf("%02d.%02d", d.Day(), d.Month())
}

func (s *ProcessorSuite) TestAddNewTask() {
	s.Run("without points", func() {
		msg := faker.Sentence()
		s.mockStore.EXPECT().CreateTask(msg, 1)
		dump := faker.Sentence()
		s.mockStore.EXPECT().CurrentSprintDump().Return(dump, nil)

		resp := s.processor.Process(msg)
		s.Require().Equal(dump, resp)
	})

	s.Run("with points at begin", func() {
		msg := faker.Sentence()
		points := rand.Int()
		s.mockStore.EXPECT().CreateTask(msg, points)
		dump := faker.Sentence()
		s.mockStore.EXPECT().CurrentSprintDump().Return(dump, nil)

		resp := s.processor.Process(strconv.Itoa(points) + " " + msg)
		s.Require().Equal(dump, resp)
	})

	s.Run("with points at end", func() {
		msg := faker.Sentence()
		points := rand.Int()
		s.mockStore.EXPECT().CreateTask(msg, points)
		dump := faker.Sentence()
		s.mockStore.EXPECT().CurrentSprintDump().Return(dump, nil)

		resp := s.processor.Process(msg + " " + strconv.Itoa(points))
		s.Require().Equal(dump, resp)
	})
}

func (s *ProcessorSuite) TestAddNewTaskError() {
	s.mockStore.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(errTestNewTask)

	resp := s.processor.Process(faker.Sentence())
	s.Require().Equal("Oops! Failed to create new task. Try later.", resp)
}

func (s *ProcessorSuite) TestDoneTask() {
	id := rand.Int()
	s.mockStore.EXPECT().DoneTask(id)
	dump := faker.Sentence()
	s.mockStore.EXPECT().CurrentSprintDump().Return(dump, nil)

	resp := s.processor.Process("/d " + strconv.Itoa(id))
	s.Require().Equal(dump, resp)
}

func (s *ProcessorSuite) TestDoneTaskInvalidID() {
	resp := s.processor.Process("/d id")
	s.Require().Equal("Task id should be a number.", resp)
}

func (s *ProcessorSuite) TestDoneTaskError() {
	s.mockStore.EXPECT().DoneTask(gomock.Any()).Return(errTestDoneTask)

	resp := s.processor.Process("/d 1")
	s.Require().Equal("Oops! Failed to done task. Try later.", resp)
}

func (s *ProcessorSuite) TestBurnTaskPoints() {
	id := rand.Int()
	burnt := rand.Int()
	s.mockStore.EXPECT().BurnTaskPoints(id, burnt)
	dump := faker.Sentence()
	s.mockStore.EXPECT().CurrentSprintDump().Return(dump, nil)

	resp := s.processor.Process("/d " + strconv.Itoa(id) + " " + strconv.Itoa(burnt))
	s.Require().Equal(dump, resp)
}

func (s *ProcessorSuite) TestBurnTaskPointsInvalidNumber() {
	resp := s.processor.Process("/d 0 burnt")
	s.Require().Equal("Burnt points should be a number.", resp)
}

func (s *ProcessorSuite) TestBurnTaskPointsError() {
	s.mockStore.EXPECT().BurnTaskPoints(gomock.Any(), gomock.Any()).Return(errTestBurnPoints)

	resp := s.processor.Process("/d 1 1")
	s.Require().Equal("Oops! Failed to burn task points. Try later.", resp)
}

func (s *ProcessorSuite) TestDoneTooMuchParams() {
	resp := s.processor.Process("/d id burnt extra")
	s.Require().Equal("Too much params for done command. Should be `/d id [burnt]`.", resp)
}

func (s *ProcessorSuite) TestCurrentTaskListError() {
	s.mockStore.EXPECT().CreateNewSprint(gomock.Any(), gomock.Any())
	s.mockStore.EXPECT().CurrentSprintDump().Return("", errTestCurrentList)
	resp := s.processor.Process("/ns 01.01 - 02.01")
	s.Require().Equal("Oops! Failed to load current task list. Try later.", resp)
}

func TestProcessorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProcessorSuite))
}
