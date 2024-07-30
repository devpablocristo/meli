package task

type MySqlRepositoryPort interface {
	CreateTask(*Task) error
	GetAllTasks() ([]Task, error)
	UpdateTask(*Task) error
	DeleteTask(string) error
}

type TaskDAOPort interface {
	Create(*TaskDaoData) error
	GetAll() ([]TaskDaoData, error)
	Update(*TaskDaoData) error
	Delete(string) error
}
