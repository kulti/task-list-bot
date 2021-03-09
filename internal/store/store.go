package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kulti/task-list-bot/internal/models"
	repo "github.com/kulti/task-list-bot/internal/repository"
)

type repository interface {
	CreateNewSprint(sprint repo.Sprint, data []byte) error
	CurrentSprint() ([]byte, error)
	UpdateCurrentSprint(data []byte) error
}

type Store struct {
	repo  repository
	tasks taskListHistoryItem
}

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

	return s.tasks.TaskList, nil
}

func (s *Store) CreateNewSprint(begin, end time.Time) error {
	tasks := taskListHistoryItem{
		Description: "Sprint is created",
		TaskList: models.TaskList{
			Title: fmt.Sprintf("%s - %s", s.timeToSprintDate(begin), s.timeToSprintDate(end)),
		},
	}

	data, err := json.Marshal(&tasks)
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}

	err = s.repo.CreateNewSprint(repo.Sprint{Begin: begin, End: end}, data)
	if err != nil {
		return fmt.Errorf("failed to create new sprint in db: %w", err)
	}

	s.tasks = tasks
	return nil
}

func (s *Store) CreateTask(text string, points int) error {
	if err := s.init(); err != nil {
		return err
	}

	histItem := s.dupHistoryItem()
	histItem.Description = fmt.Sprintf("Task '%q' is created", text)

	id := len(histItem.Tasks)
	task := models.Task{ID: id, Text: text, Points: models.Points{Total: points}}
	histItem.Tasks = append(histItem.Tasks, task)
	histItem.Points.Total += points

	return s.flush(histItem)
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
	histItem.Description = fmt.Sprintf("Mark task %q as done", task.Text)
	histItem.Points.Burnt += task.Points.Total - task.Points.Burnt
	histItem.Tasks[id].State = models.TaskStateDone
	histItem.Tasks[id].Points.Burnt = task.Points.Total

	return task.Text, s.flush(histItem)
}

func (s *Store) BurnTaskPoints(id int, burnt int) (string, error) {
	if err := s.init(); err != nil {
		return "", err
	}

	histItem := s.dupHistoryItem()

	if id >= len(histItem.Tasks) {
		return "", models.ErrTaskNotFound
	}

	histItem.Tasks = append([]models.Task{}, histItem.Tasks...)
	task := &histItem.Tasks[id]
	histItem.Description = fmt.Sprintf("Burnt %d points for task %q", burnt, task.Text)
	histItem.Points.Burnt += burnt
	task.Points.Burnt += burnt

	if task.Points.Burnt > task.Points.Total {
		histItem.Points.Total += task.Points.Burnt - task.Points.Total
		task.Points.Total = task.Points.Burnt
		histItem.Description += " ❗ burnt more than available"
	}

	if task.Points.Burnt == task.Points.Total {
		task.State = models.TaskStateDone
		histItem.Description += " and it's done!"
	}

	return task.Text, s.flush(histItem)
}

func (s *Store) timeToSprintDate(d time.Time) string {
	return fmt.Sprintf("%02d.%02d", d.Day(), d.Month())
}

func (s *Store) dupHistoryItem() taskListHistoryItem {
	return s.tasks
}

func (s *Store) init() error {
	if len(s.tasks.Description) != 0 {
		return nil
	}

	data, err := s.repo.CurrentSprint()
	if err != nil {
		return fmt.Errorf("query current sprint: %w", err)
	}

	var histItem taskListHistoryItem
	if err := json.Unmarshal(data, &histItem); err != nil {
		return fmt.Errorf("query current sprint: %w", err)
	}

	s.tasks = histItem
	return nil
}

func (s *Store) flush(histItem taskListHistoryItem) error {
	data, err := json.Marshal(&histItem)
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}

	if err := s.repo.UpdateCurrentSprint(data); err != nil {
		return fmt.Errorf("update tasks in db: %w", err)
	}

	s.tasks = histItem

	return nil
}
