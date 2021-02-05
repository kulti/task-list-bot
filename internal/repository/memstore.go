package repository

import "time"

type Store struct{}

func New() *Store {
	return &Store{}
}

func (*Store) CreateNewSprint(begin, end time.Time) error {
	return nil
}
