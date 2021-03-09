package repository

import (
	"time"
)

type Sprint struct {
	Begin, End time.Time
}

type Store struct {
	sprint Sprint
	tasks  [][]byte
}

func New() *Store {
	return &Store{}
}

func (s *Store) CreateNewSprint(sprint Sprint, data []byte) error {
	s.sprint = sprint
	s.tasks = [][]byte{data}
	return nil
}

func (s *Store) CurrentSprint() ([]byte, error) {
	return s.tasks[len(s.tasks)-1], nil
}

func (s *Store) UpdateCurrentSprint(data []byte) error {
	s.tasks = append(s.tasks, data)
	return nil
}
