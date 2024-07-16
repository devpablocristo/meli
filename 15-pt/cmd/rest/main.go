package main

import (
	"log"

	"github.com/gin-gonic/gin"

	handler "api/cmd/rest/handlers"
	core "api/internal/core"
	item "api/internal/core/item"
	mysqlsetup "api/internal/platform/mysql"
)

func main() {
	// Setup MySQL
	mysqlClient, err := mysqlsetup.NewMySQLSetup()
	if err != nil {
		log.Fatalf("Could not set up MySQL: %v", err)
	}
	defer mysqlClient.Close()

	// Initialize repository with MySQL connection
	repo := item.NewMySqlRepository(mysqlClient.DB())

	// Initialize use case
	usecase := core.NewItemUsecase(repo)

	// Initialize handlers
	handler := handler.NewHandler(usecase)

	// Setup router
	router := gin.Default()
	router.POST("/items", handler.SaveItem)
	router.GET("/items", handler.ListItems)

	// Start server
	log.Println("Server started at http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
