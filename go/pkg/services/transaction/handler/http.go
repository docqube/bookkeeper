package handler

import (
	"net/http"
	"strconv"
	"time"

	"docqube.de/bookkeeper/pkg/services/transaction"
	"docqube.de/bookkeeper/pkg/services/transaction/csv"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *transaction.Service
}

func NewHandler(router *gin.RouterGroup) *Handler {
	handler := &Handler{
		Service: &transaction.Service{},
	}

	transactionsAPI := router.Group("/transactions")
	transactionsAPI.POST("/csv", handler.ImportCSV)
	transactionsAPI.GET("/unclassified", handler.ListUnclassified)
	transactionsAPI.GET("", handler.List)

	transactionAPI := router.Group("/transaction")
	transactionAPI.GET("/:id", handler.Get)
	transactionAPI.PATCH("/:id", handler.Patch)

	return handler
}

func (h *Handler) ImportCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	csvFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, err := csv.ParseFile(csvFile, csv.INGConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.CategorizeAndImport(transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (h *Handler) List(c *gin.Context) {
	fromString := c.Query("from")
	if fromString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from is required"})
		return
	}

	toString := c.Query("to")
	if toString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to is required"})
		return
	}

	from, err := time.Parse(time.DateOnly, fromString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	to, err := time.Parse(time.DateOnly, toString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawCategory := c.Query("category")
	if rawCategory != "" {
		categoryID, err := strconv.ParseInt(rawCategory, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		transactions, err := h.Service.ListByCategoryID(from, to, categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, transactions)
		return
	}

	transactions, err := h.Service.List(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) ListUnclassified(c *gin.Context) {
	fromString := c.Query("from")
	if fromString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from is required"})
		return
	}

	toString := c.Query("to")
	if toString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to is required"})
		return
	}

	from, err := time.Parse(time.DateOnly, fromString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	to, err := time.Parse(time.DateOnly, toString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, err := h.Service.ListUnclassified(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) Patch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var patchRequest transaction.TransactionRequest
	err = c.BindJSON(&patchRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if patchRequest.CategoryID != nil {
		err = h.Service.Categorize(id, *patchRequest.CategoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	transaction, err := h.Service.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.Service.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
