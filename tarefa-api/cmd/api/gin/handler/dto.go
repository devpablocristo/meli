package handler

import (
	"github.com/devpablocristo/tarefaapi/internal/taskslist/task"
	"github.com/devpablocristo/tarefaapi/internal/taskslist/task/user"
)

type TaskDTO struct {
	ID          string  `json:"id"`
	Title       string  `json:"titulo"`
	Description string  `json:"descricao"`
	User        UserDTO `json:"user"`
	Status      string  `json:"status"`
}

type UserDTO struct {
	Name string `json:"nome"`
	Team string `json:"time"`
}

func dto2domain(dto TaskDTO) task.Task {
	user := user.User{
		Name: dto.User.Name,
		Team: dto.User.Team,
	}

	return task.Task{
		ID:          dto.ID,
		Title:       dto.Title,
		Description: dto.Description,
		User:        user,
		Status:      dto.Status,
	}
}

func domain2dto(dom *task.Task) *TaskDTO {
	user := UserDTO{
		Name: dom.User.Name,
		Team: dom.User.Team,
	}

	return &TaskDTO{
		ID:          dom.ID,
		Title:       dom.Title,
		Description: dom.Description,
		User:        user,
		Status:      dom.Status,
	}
}
