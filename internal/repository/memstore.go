package repository

type Store struct {
	data []byte
}

func New() *Store {
	return &Store{}
}

func (s *Store) InitNewSprint() error {
	s.data = nil
	return nil
}

func (s *Store) CurrentSprint() ([]byte, error) {
	return s.data, nil
}

func (s *Store) UpdateCurrentSprint(data []byte) error {
	s.data = data
	return nil
}
