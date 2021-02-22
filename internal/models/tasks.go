package models

import "fmt"

// TaskList represents a task list.
type TaskList struct {
	Title  string
	Points Points
	Tasks  []Task
}

// Task represents a task.
type Task struct {
	ID     int
	Text   string
	State  TaskState
	Points Points
}

// Points represents progress points.
type Points struct {
	Burnt int
	Total int
}

func (p Points) String() string {
	return fmt.Sprintf("(%d/%d)", p.Burnt, p.Total)
}

// TaskState reprensts a task state.
type TaskState string

// TaskState constants.
const (
	TaskStateSimple   TaskState = ""
	TaskStateTodo     TaskState = "todo"
	TaskStateDone     TaskState = "done"
	TaskStateCanceled TaskState = "canceled"
)
