package store_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

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

	s.mockRepository.EXPECT().InitNewSprint()
	s.mockRepository.EXPECT().UpdateCurrentSprint(bytesMatcher("new_sprint"))

	s.Require().NoError(s.store.CreateNewSprint(begin, end))
}

func TestStore(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StoreSuite))
}
