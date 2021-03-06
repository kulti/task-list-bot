package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kulti/task-list-bot/internal/models"
)

type repository interface {
	InitNewSprint() error
	CurrentSprint() ([]byte, error)
	UpdateCurrentSprint(data []byte) error
}

type Store struct {
	repo  repository
	tasks taskListHistory
}

type taskListHistory []*taskListHistoryItem

type taskListHistoryItem struct {
	Description string
	models.TaskList
}

func New(repo repository) *Store {
	return &Store{
		repo: repo,
	}
}

func (s *Store) CurrentSprint() (models.TaskList, error) {
	if err := s.init(); err != nil {
		return models.TaskList{}, err
	}

	return s.tasks[len(s.tasks)-1].TaskList, nil
}

func (s *Store) CreateNewSprint(begin, end time.Time) error {
	s.tasks = []*taskListHistoryItem{{
		Description: "Sprint is created",
		TaskList: models.TaskList{
			Title: fmt.Sprintf("%s - %s", s.timeToSprintDate(begin), s.timeToSprintDate(end)),
		},
	}}

	_ = s.repo.InitNewSprint()
	if err := s.flush(); err != nil {
		s.tasks = nil
		return fmt.Errorf("flush task list into db: %w", err)
	}

	return nil
}

func (s *Store) CreateTask(text string, points int) error {
	if err := s.init(); err != nil {
		return err
	}

	histItem := s.dupHistoryItem()

	id := len(histItem.Tasks)
	task := models.Task{ID: id, Text: text, Points: models.Points{Total: points}}
	histItem.Tasks = append(histItem.Tasks, task)
	histItem.Points.Total += points

	return s.putHistoryItem(histItem)
}

func (s *Store) DoneTask(id int) (string, error) {
	if err := s.init(); err != nil {
		return "", err
	}

	histItem := s.dupHistoryItem()

	if id >= len(histItem.Tasks) {
		return "", models.ErrTaskNotFound
	}

	histItem.Tasks = append([]models.Task{}, histItem.Tasks...)
	task := histItem.Tasks[id]
	histItem.Points.Burnt += task.Points.Total - task.Points.Burnt
	histItem.Tasks[id].State = models.TaskStateDone
	histItem.Tasks[id].Points.Burnt = task.Points.Total

	return task.Text, s.putHistoryItem(histItem)
}

func (s *Store) timeToSprintDate(d time.Time) string {
	return fmt.Sprintf("%02d.%02d", d.Day(), d.Month())
}

func (s *Store) dupHistoryItem() *taskListHistoryItem {
	histItem := *s.tasks[len(s.tasks)-1]
	return &histItem
}

func (s *Store) putHistoryItem(histItem *taskListHistoryItem) error {
	s.tasks = append(s.tasks, histItem)
	if err := s.flush(); err != nil {
		s.tasks = s.tasks[:len(s.tasks)-1]
		return fmt.Errorf("flush task list into db: %w", err)
	}
	return nil
}

func (s *Store) init() error {
	if len(s.tasks) != 0 {
		return nil
	}

	data, err := s.repo.CurrentSprint()
	if err != nil {
		return fmt.Errorf("query current sprint: %w", err)
	}

	var tasksHistory taskListHistory
	if err := json.Unmarshal(data, &tasksHistory); err != nil {
		return fmt.Errorf("query current sprint: %w", err)
	}

	s.tasks = tasksHistory
	return nil
}

func (s *Store) flush() error {
	data, err := json.Marshal(&s.tasks)
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}

	return s.repo.UpdateCurrentSprint(data)
}
