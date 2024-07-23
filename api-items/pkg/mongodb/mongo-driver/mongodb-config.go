package mongodbdriver

import "fmt"

// MongoDBClientConfig contém a configuração necessária para se conectar a um banco de dados MongoDB
type MongoDBClientConfig struct {
	User     string // Usuário do banco de dados
	Password string // Senha do usuário
	Host     string // Host onde o banco de dados está localizado
	Port     string // Porta na qual o banco de dados está ouvindo
	Database string // Nome do banco de dados
}

// dsn gera o Data Source Name (DSN) a partir da configuração fornecida
func (config MongoDBClientConfig) dsn() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, config.Database)

	//return "mongodb://root:rootpassword@mongodb:27017"
}
