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

	"github.com/kulti/task-list-bot/internal/models"
	"github.com/kulti/task-list-bot/internal/processor"
)

const noSprintTemplate = "no sprint\n"

var (
	errTestNewSprint   = errors.New("test: failed to store sprint")
	errTestNewTask     = errors.New("test: failed to create task")
	errTestDoneTask    = errors.New("test: failed to done task")
	errTestCurrentList = errors.New("test: failed to list tasks")
)

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
	s.mockRepository.EXPECT().CurrentSprint()

	sprintHeader := fmt.Sprintf("%s - %s", s.timeToSprintDate(begin), s.timeToSprintDate(end))
	resp := s.processor.Process("/ns " + sprintHeader)
	s.Require().Equal(noSprintTemplate, resp)
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

func (s *ProcessorSuite) TestAddNewTask() {
	s.Run("without points", func() {
		msg := faker.Sentence()
		s.mockRepository.EXPECT().CreateTask(msg, 1)
		s.mockRepository.EXPECT().CurrentSprint()

		resp := s.processor.Process(msg)
		s.Require().Equal(noSprintTemplate, resp)
	})

	s.Run("with points at begin", func() {
		msg := faker.Sentence()
		points := rand.Int()
		s.mockRepository.EXPECT().CreateTask(msg, points)
		s.mockRepository.EXPECT().CurrentSprint()

		resp := s.processor.Process(strconv.Itoa(points) + " " + msg)
		s.Require().Equal(noSprintTemplate, resp)
	})

	s.Run("with points at end", func() {
		msg := faker.Sentence()
		points := rand.Int()
		s.mockRepository.EXPECT().CreateTask(msg, points)
		s.mockRepository.EXPECT().CurrentSprint()

		resp := s.processor.Process(msg + " " + strconv.Itoa(points))
		s.Require().Equal(noSprintTemplate, resp)
	})
}

func (s *ProcessorSuite) TestAddNewTaskError() {
	s.mockRepository.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(errTestNewTask)

	resp := s.processor.Process(faker.Sentence())
	s.Require().Equal("Oops! Failed to create new task. Try later.", resp)
}

func (s *ProcessorSuite) TestDoneTask() {
	id := rand.Int()
	task := faker.Sentence()
	s.mockRepository.EXPECT().DoneTask(id).Return(task, nil)
	s.mockRepository.EXPECT().CurrentSprint()

	resp := s.processor.Process("/d " + strconv.Itoa(id))
	s.Require().Equal(fmt.Sprintf("The task %q is marked as done.\n\n%s", task, noSprintTemplate), resp)
}

func (s *ProcessorSuite) TestDoneTaskInvalidID() {
	resp := s.processor.Process("/d id")
	s.Require().Equal("Task id should be a number.", resp)
}

func (s *ProcessorSuite) TestDoneTaskError() {
	s.mockRepository.EXPECT().DoneTask(gomock.Any()).Return("", errTestDoneTask)

	resp := s.processor.Process("/d 1")
	s.Require().Equal("Oops! Failed to done task. Try later.", resp)
}

func (s *ProcessorSuite) TestCurrentTaskList() {
	taskList := models.TaskList{
		Title: "test title",
		Tasks: []models.Task{
			{ID: 0, Text: "task 7", Points: 7},
			{ID: 1, Text: "burnt task", Points: 3, Burnt: 2},
		},
	}
	s.mockRepository.EXPECT().CreateNewSprint(gomock.Any(), gomock.Any())
	s.mockRepository.EXPECT().CurrentSprint().Return(taskList, nil)

	resp := s.processor.Process("/ns 01.01 - 02.01")
	s.Require().Equal(`test title
0 (0/7) task 7
1 (2/3) burnt task
`, resp)
}

func (s *ProcessorSuite) TestCurrentTaskListError() {
	s.mockRepository.EXPECT().CreateNewSprint(gomock.Any(), gomock.Any())
	s.mockRepository.EXPECT().CurrentSprint().Return(models.TaskList{}, errTestCurrentList)
	resp := s.processor.Process("/ns 01.01 - 02.01")
	s.Require().Equal("Oops! Failed to load current task list. Try later.", resp)
}

func TestProcessorSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProcessorSuite))
}
