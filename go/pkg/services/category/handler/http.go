package handler

import (
	"database/sql"
	"net/http"

	"docqube.de/bookkeeper/pkg/services/category"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *category.Service
}

func NewHandler(router *gin.RouterGroup, db *sql.DB) *Handler {
	handler := &Handler{
		service: category.NewService(db),
	}

	categoriesAPI := router.Group("/categories")
	categoriesAPI.GET("", handler.List)

	return handler
}

func (h *Handler) List(c *gin.Context) {
	categories, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
