package mysqlsetup

import (
	gosqldriver "api/pkg/mysql/go-sql-driver"
)

// NewMySQLSetup configura e retorna um novo cliente MySQL
func NewMySQLSetup() (*gosqldriver.MySQLClient, error) {
	config := gosqldriver.MySQLClientConfig{
		User:     "api_user",
		Password: "api_password",
		Host:     "mysql",
		Port:     "3306",
		Database: "inventory",
	}
	return gosqldriver.NewMySQLClient(config)
}
