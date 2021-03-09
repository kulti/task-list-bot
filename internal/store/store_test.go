package store_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	repo "github.com/kulti/task-list-bot/internal/repository"
	"github.com/kulti/task-list-bot/internal/store"
)

type StoreSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	mockRepository *MockRepository
	store          *store.Store
}

func (s *StoreSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockRepository = NewMockRepository(s.mockCtrl)
	s.store = store.New(s.mockRepository)
}

func (s *StoreSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *StoreSuite) TestCreateNewSprint() {
	const timeDay = 24 * time.Hour
	begin := time.Unix(13*int64(timeDay.Seconds()), 0)
	end := begin.Add(7 * timeDay)

	s.mockRepository.EXPECT().CreateNewSprint(repo.Sprint{Begin: begin, End: end}, bytesMatcher("new_sprint"))

	s.Require().NoError(s.store.CreateNewSprint(begin, end))
}

func (s *StoreSuite) TestCreateTask() {
	sprintData, err := ioutil.ReadFile(filepath.Join("testdata", "new_sprint"))
	s.Require().NoError(err)
	s.mockRepository.EXPECT().CurrentSprint().Return(sprintData, nil)

	s.mockRepository.EXPECT().UpdateCurrentSprint(bytesMatcher("new_task"))

	s.Require().NoError(s.store.CreateTask("new task", 10))
}

func (s *StoreSuite) TestDoneTask() {
	sprintData, err := ioutil.ReadFile(filepath.Join("testdata", "new_task"))
	s.Require().NoError(err)
	s.mockRepository.EXPECT().CurrentSprint().Return(sprintData, nil)

	s.mockRepository.EXPECT().UpdateCurrentSprint(bytesMatcher("done_task"))

	task, err := s.store.DoneTask(0)
	s.Require().NoError(err)
	s.Require().Equal("new task", task)
}

func (s *StoreSuite) TestBurnPoints() {
	sprintData, err := ioutil.ReadFile(filepath.Join("testdata", "burn_init"))
	s.Require().NoError(err)
	s.mockRepository.EXPECT().CurrentSprint().Return(sprintData, nil)

	s.mockRepository.EXPECT().UpdateCurrentSprint(bytesMatcher("burn_some_task_0"))
	task0, err := s.store.BurnTaskPoints(0, 1)
	s.Require().NoError(err)
	s.Require().Equal("task 0", task0)

	s.mockRepository.EXPECT().UpdateCurrentSprint(bytesMatcher("burn_overflow_task_1"))
	task1, err := s.store.BurnTaskPoints(1, 2)
	s.Require().NoError(err)
	s.Require().Equal("task 1", task1)

	s.mockRepository.EXPECT().UpdateCurrentSprint(bytesMatcher("burn_done_task_0"))
	task0, err = s.store.BurnTaskPoints(0, 3)
	s.Require().NoError(err)
	s.Require().Equal("task 0", task0)
}

func TestStore(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StoreSuite))
}
