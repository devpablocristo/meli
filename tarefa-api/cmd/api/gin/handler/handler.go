package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/devpablocristo/tarefaapi/internal/taskslist"
)

type Handler struct {
	usecase taskslist.TaskUsecasePort
}

func NewHandler(u taskslist.TaskUsecasePort) *Handler {
	return &Handler{
		usecase: u,
	}
}

func (h *Handler) CreateTask(c *gin.Context) {
	var taskDTO TaskDTO
	// Checks for error decoding JSON
	if err := c.BindJSON(&taskDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskDomain := dto2domain(taskDTO)
	err := h.usecase.CreateTask(&taskDomain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain2dto(&taskDomain))
}

func (h *Handler) UpdateTask(c *gin.Context) {
	var taskDTO TaskDTO

	if err := c.BindJSON(&taskDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskDTO.ID = c.Param("ID")

	taskDomain := dto2domain(taskDTO)
	err := h.usecase.UpdateTask(&taskDomain)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain2dto(&taskDomain))
}

func (h *Handler) GetAllTask(c *gin.Context) {
	task, err := h.usecase.GetAllTask()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	ID := c.Param("ID")
	err := h.usecase.DeleteTask(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "deleted task successfully")
}
