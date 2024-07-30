package gin

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/devpablocristo/tarefaapi/cmd/api/gin/handler"
	"github.com/devpablocristo/tarefaapi/internal/platform/env"
)

func NewHTTPServer(h handler.Handler) error {
	port := env.GetEnv("PORT", ":8080")

	r := gin.Default()

	setupRoutes(r, &h)

	log.Printf("Server listening on port %s", port)
	return r.Run(port)
}

func setupRoutes(r *gin.Engine, h *handler.Handler) {
	basePath := "/api"
	router := r.Group(basePath)

	router.POST("/tarefas", h.CreateTask)
	router.PUT("/tarefas/:ID", h.UpdateTask)
	router.GET("/tarefas/:status", h.GetAllTask)
	router.DELETE("/tarefas/:ID", h.DeleteTask)
}
