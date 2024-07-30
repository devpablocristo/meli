package task

import (
	"time"

	"github.com/devpablocristo/tarefaapi/internal/taskslist/task/user"
)

type TaskDaoData struct {
	ID          string      `db:"ID" bson:"ID"`
	Title       string      `db:"title" bson:"title"`
	Description string      `db:"description" bson:"description"`
	User        UserDaoData `db:"user" bson:"user"`
	Status      string      `db:"status" bson:"status"`
	CreatedAt   time.Time   `db:"created_at" bson:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at" bson:"updated_at"`
	Deleted     bool        `db:"deleted" bson:"deleted"`
}

type UserDaoData struct {
	Name string `db:"name" bson:"name"`
	Team string `db:"team" bson:"team"`
}

func Domain2dao(dom *Task) *TaskDaoData {
	return &TaskDaoData{
		ID:          dom.ID,
		Title:       dom.Title,
		Description: dom.Description,
		User: UserDaoData{
			Name: dom.User.Name,
			Team: dom.User.Team,
		},
		Status: dom.Status,
	}
}

func Dao2domain(dao *TaskDaoData) *Task {
	return &Task{
		ID:          dao.ID,
		Title:       dao.Title,
		Description: dao.Description,
		User: user.User{
			Name: dao.User.Name,
			Team: dao.User.Team,
		},
		Status: dao.Status,
	}
}

func Dao2domainList(daos []TaskDaoData) []Task {
	var task []Task
	for _, dao := range daos {
		task = append(task, *Dao2domain(&dao))
	}
	return task
}
