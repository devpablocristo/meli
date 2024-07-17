package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api/internal/core"
	"api/internal/core/item"
	"api/pkg/config"
)

// handler es el manejador para las solicitudes HTTP relacionadas con los elementos
type handler struct {
	core core.ItemUsecasePort // Caso de uso de elementos
}

// NewHandler crea una nueva instancia de handler
func NewHandler(u core.ItemUsecasePort) *handler {
	return &handler{
		core: u,
	}
}

// SaveItem maneja la solicitud para guardar un nuevo elemento
func (h *handler) SaveItem(c *gin.Context) {
	var it item.Item
	err := c.BindJSON(&it)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.core.SaveItem(it); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "item saved successfully")
}

// ListItems maneja la solicitud para listar todos los elementos
func (h *handler) ListItems(c *gin.Context) {
	its, err := h.core.ListItems()
	if err != nil {
		if err == config.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, its)
}
