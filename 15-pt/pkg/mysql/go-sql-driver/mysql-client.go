package gosqldriver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLClient representa um cliente para interagir com um banco de dados MySQL
type MySQLClient struct {
	config MySQLClientConfig // Configuração do cliente MySQL
	db     *sql.DB           // Conexão com o banco de dados
}

// NewMySQLClient cria uma nova instância de MySQLClient e estabelece a conexão com o banco de dados
func NewMySQLClient(config MySQLClientConfig) (*MySQLClient, error) {
	client := &MySQLClient{config: config}
	err := client.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MySQLClient: %v", err)
	}
	return client, nil
}

// connect estabelece a conexão com o banco de dados MySQL utilizando a configuração fornecida
func (client *MySQLClient) connect() error {
	dsn := client.config.dsn()
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}
	client.db = conn
	return nil
}

// Close fecha a conexão com o banco de dados MySQL
func (client *MySQLClient) Close() {
	if client.db != nil {
		client.db.Close()
	}
}

// DB retorna a conexão com o banco de dados MySQL
func (client *MySQLClient) DB() *sql.DB {
	return client.db
}
