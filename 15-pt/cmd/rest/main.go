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
		log.Fatalf("No se pudo configurar MySQL: %v", err)
	}
	defer mysqlClient.Close()

	// Inicializar repositorios
	mapRepo := item.NewRepository()
	mysqlRepo := item.NewMySqlRepository(mysqlClient.DB())

	// Inicializar caso de uso con ambos repositorios
	usecase := core.NewItemUsecase(mysqlRepo, mapRepo)

	// Inicializar handlers
	handler := handler.NewHandler(usecase)

	// Configurar enrutador
	router := gin.Default()
	router.POST("/items", handler.SaveItem)
	router.GET("/items", handler.ListItems)

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
