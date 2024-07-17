package gosqldriver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLClient representa un cliente para interactuar con una base de datos MySQL
type MySQLClient struct {
	config MySQLClientConfig // Configuración del cliente MySQL
	db     *sql.DB           // Conexión a la base de datos
}

// NewMySQLClient crea una nueva instancia de MySQLClient y establece la conexión a la base de datos
func NewMySQLClient(config MySQLClientConfig) (*MySQLClient, error) {
	client := &MySQLClient{config: config}
	err := client.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MySQLClient: %v", err)
	}
	return client, nil
}

// connect establece la conexión a la base de datos MySQL utilizando la configuración proporcionada
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

// Close cierra la conexión a la base de datos MySQL
func (client *MySQLClient) Close() {
	if client.db != nil {
		client.db.Close()
	}
}

// DB devuelve la conexión a la base de datos MySQL
func (client *MySQLClient) DB() *sql.DB {
	return client.db
}
