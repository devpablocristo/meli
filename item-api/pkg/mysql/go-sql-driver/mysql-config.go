package gosqldriver

import (
	"fmt"
)

// MySQLClientConfig contém a configuração necessária para se conectar a um banco de dados MySQL
type MySQLClientConfig struct {
	User     string // Usuário do banco de dados
	Password string // Senha do usuário
	Host     string // Host onde o banco de dados está localizado
	Port     string // Porta na qual o banco de dados está ouvindo
	Database string // Nome do banco de dados
}

// dsn gera o Data Source Name (DSN) a partir da configuração fornecida
func (config MySQLClientConfig) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Database)
}
