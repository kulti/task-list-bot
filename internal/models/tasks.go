package models

// TaskList represents a task list.
type TaskList struct {
	Title string
	Tasks []Task
}

// Task represents a task.
type Task struct {
	ID     int
	Text   string
	State  TaskState
	Points int
	Burnt  int
}

// TaskState reprensts a task state.
type TaskState string

// TaskState constants.
const (
	TaskStateSimple    TaskState = ""
	TaskStateTodo      TaskState = "todo"
	TaskStateCompleted TaskState = "done"
	TaskStateCanceled  TaskState = "canceled"
)
