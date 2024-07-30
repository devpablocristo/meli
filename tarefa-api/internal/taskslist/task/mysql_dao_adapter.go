package task

import (
	"github.com/jmoiron/sqlx"
)

type MySqlDao struct {
	db *sqlx.DB
}

func NewMySqlDao(db *sqlx.DB) TaskDAOPort {
	return &MySqlDao{
		db: db,
	}
}

func (dao *MySqlDao) Create(d *TaskDaoData) error {
	query := `INSERT INTO task (ID, title, description, name, team, status, created_at, updated_at, deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := dao.db.Exec(query, d.ID, d.Title, d.Description, d.User.Name, d.User.Team, d.Status, d.CreatedAt, d.UpdatedAt, d.Deleted)
	if err != nil {
		return err
	}
	return err
}

func (dao *MySqlDao) GetAll() ([]TaskDaoData, error) {
	var d []TaskDaoData
	query := "SELECT * FROM task WHERE deleted = FALSE"
	err := dao.db.Select(&d, query)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (dao *MySqlDao) Update(d *TaskDaoData) error {
	query := `UPDATE task SET title = ?, description = ?, name = ?, team = ?, status = ?, updated_at = ?, deleted = ? WHERE ID = ?`
	_, err := dao.db.Exec(query, d.Title, d.Description, d.User.Name, d.User.Team, d.Status, d.UpdatedAt, d.Deleted, d.ID)
	if err != nil {
		return err
	}
	return err
}

func (dao *MySqlDao) Delete(ID string) error {
	query := `UPDATE task SET deleted = TRUE WHERE ID = ?`
	_, err := dao.db.Exec(query, ID)
	if err != nil {
		return err
	}
	return err
}
