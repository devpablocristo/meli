package mysql

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/devpablocristo/tarefaapi/internal/platform/env"
)

var (
	DBHost = env.GetEnv("MYSQL_HOST", "mysql-local")
	DBPort = env.GetEnv("MYSQL_CONT_PORT", "3306")
	DBName = env.GetEnv("MYSQL_DATABASE", "tarefaapi")
	DBUser = env.GetEnv("MYSQL_USERNAME", "root")
	DBPass = env.GetEnv("MYSQL_USER_PASSWORD", "root")
)

var (
	db   *sqlx.DB
	once sync.Once // Used to ensure that the DB initialization happens only once
)

// GetConnectionDB establishes and returns a singleton database connection.
// This implements the Singleton pattern to ensure only one instance of the connection.
func GetConnectionDB() (*sqlx.DB, error) {
	var err error
	once.Do(func() { // Ensures that the connection is only created once
		dsn := dbConnectionURL() // Builds the DSN for the connection
		db, err = sqlx.Connect("mysql", dsn)
		if err != nil {
			log.Printf("DB Connection Error: %s", err)
			return
		}

		// Optional: Invoke autoMigrate here if you want to apply migrations
		// automatically every time the database connection is established.
		// However, consider manual control for production environments.
		// if err := autoMigrate(); err != nil {
		//     log.Printf("Failed to auto-migrate database schema: %v", err)
		//     // Consider if you want a migration failure to stop execution
		// }
	})
	return db, err
}

// autoMigrate automatically applies database migrations.
// This function is crucial for evolving your database schema safely and efficiently.
// It uses the golang-migrate library to apply migrations stored in a specific directory.
// func autoMigrate() error {
// 	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
// 	if err != nil {
// 		log.Printf("Migration Driver Error: %s", err)
// 		return fmt.Errorf("migration Driver Error: %w", err)
// 	}

// 	m, err := migrate.NewWithDatabaseInstance(
// 		"file://path/to/your/migrations", // Adjust this path to your migration files
// 		"mysql", driver,
// 	)
// 	if err != nil {
// 		log.Printf("Migration Instance Error: %s", err)
// 		return fmt.Errorf("migration Instance Error: %w", err)
// 	}

// 	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 		log.Printf("Migration Error: %s", err)
// 		return fmt.Errorf("migration Error: %w", err)
// 	}

// 	return nil
// }

// dbConnectionURL builds the MySQL connection string.
// This function constructs the DSN (Data Source Name) for connecting to MySQL.
// It fetches configuration from environment variables, providing defaults if not set.
func dbConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&multiStatements=true",
		DBUser, DBPass, DBHost, DBPort, DBName)
}
