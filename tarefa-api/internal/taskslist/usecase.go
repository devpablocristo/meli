package taskslist

import (
	"fmt"

	"github.com/devpablocristo/tarefaapi/internal/taskslist/task"
)

type TaskUsecase struct {
	Repo task.MySqlRepositoryPort
}

func NewTaskUsecase(r task.MySqlRepositoryPort) TaskUsecasePort {
	return &TaskUsecase{
		Repo: r,
	}
}

func (u *TaskUsecase) CreateTask(t *task.Task) error {
	err := u.Repo.CreateTask(t)
	if err != nil {
		return fmt.Errorf("error saving Repo task: %w", err)
	}
	return nil
}

func (u *TaskUsecase) UpdateTask(t *task.Task) error {
	err := u.Repo.UpdateTask(t)
	if err != nil {
		return fmt.Errorf("error updating task REPO: %w", err)
	}
	return nil
}

func (u *TaskUsecase) GetAllTask() ([]task.Task, error) {
	task, err := u.Repo.GetAllTasks()
	if err != nil {
		return task, fmt.Errorf("error getting task from Repo: %w", err)
	}
	return task, nil
}

func (u *TaskUsecase) DeleteTask(ID string) error {
	err := u.Repo.DeleteTask(ID)
	if err != nil {
		return fmt.Errorf("error updating task REPO: %w", err)
	}
	return nil
}
