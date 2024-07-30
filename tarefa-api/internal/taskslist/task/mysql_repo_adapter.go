package task

import (
	"time"

	"github.com/google/uuid"
)

type MySqlRepository struct {
	taskDAO TaskDAOPort
}

func NewMysqlRepository(d TaskDAOPort) MySqlRepositoryPort {
	return &MySqlRepository{
		taskDAO: d,
	}
}

func (r *MySqlRepository) CreateTask(t *Task) error {
	d := Domain2dao(t)
	d.ID = uuid.New().String()
	d.Status = "TODO"
	d.CreatedAt = time.Now()
	d.UpdatedAt = d.CreatedAt
	return r.taskDAO.Create(d)
}

func (r *MySqlRepository) GetAllTasks() ([]Task, error) {
	d, err := r.taskDAO.GetAll()
	if err != nil {
		return nil, err
	}
	return Dao2domainList(d), nil
}

func (r *MySqlRepository) UpdateTask(t *Task) error {
	d := Domain2dao(t)
	d.UpdatedAt = time.Now()
	return r.taskDAO.Update(d)
}

func (r *MySqlRepository) DeleteTask(ID string) error {
	return r.taskDAO.Delete(ID)
}
