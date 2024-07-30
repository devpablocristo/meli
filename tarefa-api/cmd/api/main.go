package main

import (
	"fmt"
	"log"

	// "github.com/devpablocristo/tarefaapi/internal/infra/mongodb"

	"github.com/devpablocristo/tarefaapi/cmd/api/gin"
	"github.com/devpablocristo/tarefaapi/cmd/api/gin/handler"
	"github.com/devpablocristo/tarefaapi/internal/platform/env"
	"github.com/devpablocristo/tarefaapi/internal/platform/mysql"
	"github.com/devpablocristo/tarefaapi/internal/taskslist"
	"github.com/devpablocristo/tarefaapi/internal/taskslist/task"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
	if err := env.LoadEnv(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	dbConn, err := setupDatabase()
	if err != nil {
		return fmt.Errorf("error setting up database: %w", err)
	}

	u := taskslist.NewTaskUsecase(dbConn)
	h := handler.NewHandler(u)

	if err := gin.NewHTTPServer(*h); err != nil {
		return fmt.Errorf("error starting HTTP server: %w", err)
	}

	return nil
}

func setupDatabase() (task.MySqlRepositoryPort, error) {
	// //MONGO DB
	// dbConn, err := mongodb.NewMongoDBConnection()
	// if err != nil {
	// 	log.Fatal("Error opening Mongo DB connection.", err)
	// }
	// r := repository.NewMongoRepo(dbConn)

	// MYSQL DB
	dbConn, err := mysql.GetConnectionDB()
	if err != nil {
		return nil, fmt.Errorf("error opening MySQL connection: %w", err)
	}
	d := task.NewMySqlDao(dbConn)
	r := task.NewMysqlRepository(d)

	// // MAP REPO
	// r := repository.NewMapRepo()

	// // SLICE REPO
	// r := repository.NewSliceRepo()

	return r, nil
}
