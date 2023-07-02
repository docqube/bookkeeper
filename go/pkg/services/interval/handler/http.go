package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"docqube.de/bookkeeper/pkg/services/interval"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *interval.Service
}

func NewHandler(router *gin.RouterGroup, db *sql.DB) *Handler {
	handler := &Handler{
		Service: interval.NewService(db),
	}

	intervalAPI := router.Group("/interval")
	intervalAPI.GET("/fiscal-month", handler.GetFiscalMonth)

	return handler
}

func (h *Handler) GetFiscalMonth(c *gin.Context) {
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incomeCategoryID, err := strconv.ParseInt(c.Query("income_category_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start, end, err := h.Service.GetFiscalMonthWithIncomeCategoryID(month, year, incomeCategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interval.FiscalMonth{
		Month: month,
		Year:  year,
		Start: *start,
		End:   *end,
	})
}
