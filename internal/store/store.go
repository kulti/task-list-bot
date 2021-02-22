package store

import (
	"time"

	"github.com/kulti/task-list-bot/internal/models"
)

type repository interface {
	CurrentSprint() (models.TaskList, error)
	CreateNewSprint(begin, end time.Time) error
	CreateTask(text string, points int) error
	DoneTask(id int) (string, error)
}

type Store struct {
	repo repository
}

func New(repo repository) *Store {
	return &Store{
		repo: repo,
	}
}

func (s *Store) CurrentSprint() (models.TaskList, error) {
	return s.repo.CurrentSprint()
}

func (s *Store) CreateNewSprint(begin, end time.Time) error {
	return s.repo.CreateNewSprint(begin, end)
}

func (s *Store) CreateTask(text string, points int) error {
	return s.repo.CreateTask(text, points)
}

func (s *Store) DoneTask(id int) (string, error) {
	return s.repo.DoneTask(id)
}
