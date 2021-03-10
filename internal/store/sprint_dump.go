package store

import (
	"strconv"
	"strings"
)

func (s *Store) CurrentSprintDump() (string, error) {
	if err := s.init(); err != nil {
		return "", err
	}

	if s.tasks.Title == "" {
		return "no sprint", nil
	}

	b := strings.Builder{}

	b.WriteString(s.tasks.Description)
	b.WriteByte('\n')
	b.WriteByte('\n')

	b.WriteString(s.tasks.Title)
	b.WriteByte('\n')

	b.WriteString("Total ")
	b.WriteString(s.tasks.Points.String())
	b.WriteByte('\n')
	b.WriteByte('\n')

	for _, t := range s.tasks.Tasks {
		b.WriteString(strconv.Itoa(t.ID))
		b.WriteByte(' ')
		b.WriteString(t.Points.String())
		b.WriteByte(' ')
		b.WriteString(t.Text)
		b.WriteByte('\n')
	}
	return b.String(), nil
}
