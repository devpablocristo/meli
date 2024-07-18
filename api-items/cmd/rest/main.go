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
	// Configurar MySQL
	mysqlClient, err := mysqlsetup.NewMySQLSetup()
	if err != nil {
		log.Fatalf("Não foi possível configurar o MySQL: %v", err)
	}
	defer mysqlClient.Close()

	// Inicializar repositórios
	mapRepo := item.NewRepository()
	mysqlRepo := item.NewMySqlRepository(mysqlClient.DB())

	// Inicializar caso de uso com ambos repositórios
	usecase := core.NewItemUsecase(mysqlRepo, mapRepo)

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
