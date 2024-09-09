package mysqlrepo

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	dbHost = "DB_HOST"
	dbUser = "DB_USER"
	dbPass = "DB_PASS"
	dbName = "DB_NAME"
)

var db *sqlx.DB //nolint:gochecknoglobals

func GetConnectionDB() (*sqlx.DB, error) {
	var err error

	if db == nil {
		conn := dbConnectionURL()
		db, err = sqlx.Connect("mysql", conn)
		if err != nil {
			fmt.Println(conn)
			fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
			return nil, fmt.Errorf("### DB ERROR: %w", err)
		}
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sqlx.DB) error {
	var itemsSchema = `
	CREATE TABLE IF NOT EXISTS items (
		id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
		code varchar(191) DEFAULT NULL,
		title text,
		description text,
		price decimal(10, 2) NOT NULL,
		stock bigint(20) DEFAULT NULL,
		available boolean NOT NULL,
		categories text NOT NULL,
		created_at datetime(3) DEFAULT NULL,
		updated_at datetime(3) DEFAULT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY code (code)
	  );`

	_, err := db.Exec(itemsSchema)
	if err != nil {
		fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
		return fmt.Errorf("### MIGRATION ERROR: %w", err)
	}

	return nil
}

func dbConnectionURL() string {
	if os.Getenv("GO_ENVIRONMENT") == "" {
		if err := godotenv.Load(".env"); err != nil {
			panic(err.Error())
		}
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", os.Getenv(dbUser), os.Getenv(dbPass), os.Getenv(dbHost), os.Getenv(dbName))
}
