package repository

import (
	"fmt"
	"time"

	"github.com/kulti/task-list-bot/internal/models"
)

type Store struct {
	tasks models.TaskList
}

func New() *Store {
	return &Store{}
}

func (s *Store) CreateNewSprint(begin, end time.Time) error {
	s.tasks.Title = fmt.Sprintf("%s - %s", s.timeToSprintDate(begin), s.timeToSprintDate(end))
	s.tasks.Tasks = nil
	return nil
}

func (s *Store) CreateTask(text string, points int) error {
	id := len(s.tasks.Tasks)
	s.tasks.Tasks = append(s.tasks.Tasks, models.Task{ID: id, Text: text, Points: points})
	return nil
}

func (s *Store) DoneTask(id int) (string, error) {
	if id >= len(s.tasks.Tasks) {
		return "", models.ErrTaskNotFound
	}

	s.tasks.Tasks[id].State = models.TaskStateDone

	return s.tasks.Tasks[id].Text, nil
}

func (s *Store) CurrentSprint() (models.TaskList, error) {
	return s.tasks, nil
}

func (s *Store) timeToSprintDate(d time.Time) string {
	return fmt.Sprintf("%02d.%02d", d.Day(), d.Month())
}
