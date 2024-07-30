package task

import (
	"github.com/devpablocristo/tarefaapi/internal/taskslist/task/user"
)

type TaskStatus string

// These are the possible statuses for a Task.
const (
	TaskStatusTodo  TaskStatus = "TODO"
	TaskStatusDoing TaskStatus = "DOING"
	TaskStatusDone  TaskStatus = "DONE"
)

type Task struct {
	ID          string
	Title       string
	Description string
	User        user.User
	Status      string
}
