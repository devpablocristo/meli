package main

import (
	"log"

	"github.com/gin-gonic/gin"

	handler "api/cmd/rest/handlers"
	core "api/internal/core"
	item "api/internal/core/item"
	mongodbsetup "api/internal/platform/mongodb"
	mysqlsetup "api/internal/platform/mysql"
)

func main() {
	// Configurar MySQL
	mysqlClient, err := mysqlsetup.NewMySQLSetup()
	if err != nil {
		log.Fatalf("não foi possível configurar o MySQL: %v", err)
	}
	defer mysqlClient.Close()

	mongoDBClient, err := mongodbsetup.NewMongoDBSetup()
	if err != nil {
		log.Fatalf("não foi possível configurar o MongoDB: %v", err)
	}
	defer mongoDBClient.Close()

	// Inicializar repositórios
	inMemoryRepo := item.NewInMemoryRepository()
	mysqlRepo := item.NewMySqlRepository(mysqlClient.DB())
	mongoDbRepo := item.NewMongoDbRepository(mongoDBClient.DB())

	_ = inMemoryRepo

	// Inicializar caso de uso com ambos repositórios
	usecase := core.NewItemUsecase(mysqlRepo, mongoDbRepo)

	// Inicializar handlers
	handler := handler.NewHandler(usecase)

	// Configurar roteador
	router := gin.Default()
	router.POST("/items", handler.SaveItem)
	router.GET("/items", handler.ListItems)

	// Iniciar servidor
	log.Println("Servidor iniciado em http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
