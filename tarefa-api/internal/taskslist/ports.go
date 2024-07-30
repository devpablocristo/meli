package taskslist

import "github.com/devpablocristo/tarefaapi/internal/taskslist/task"

type TaskUsecasePort interface {
	CreateTask(*task.Task) error
	UpdateTask(*task.Task) error
	GetAllTask() ([]task.Task, error)
	DeleteTask(string) error
}
