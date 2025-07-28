package http

import (
	"arise_tech_assetment/internal/application/commands"
	"arise_tech_assetment/internal/application/queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

type TransactionHandler struct {
}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction (deposit, withdraw, or transfer)
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body commands.CreateTransactionCommand true "Transaction creation data"
// @Success 201 {object} commands.CreateTransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var cmd commands.CreateTransactionCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := mediatr.Send[*commands.CreateTransactionCommand, *commands.CreateTransactionResponse](c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetTransaction godoc
// @Summary Get transaction by ID
// @Description Get a single transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} queries.GetTransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	query := &queries.GetTransactionQuery{ID: id}
	result, err := mediatr.Send[*queries.GetTransactionQuery, *queries.GetTransactionResponse](c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTransactions godoc
// @Summary Get all transactions
// @Description Get a paginated list of all transactions
// @Tags transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} queries.GetTransactionsResponse
// @Failure 500 {object} map[string]string
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := &queries.GetTransactionsQuery{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := mediatr.Send[*queries.GetTransactionsQuery, *queries.GetTransactionsResponse](c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAccountTransactions godoc
// @Summary Get transactions for an account
// @Description Get a paginated list of transactions for a specific account
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} queries.GetAccountTransactionsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id}/transactions [get]
func (h *TransactionHandler) GetAccountTransactions(c *gin.Context) {
	accountIDParam := c.Param("id")
	accountID, err := uuid.Parse(accountIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := &queries.GetAccountTransactionsQuery{
		AccountID: accountID,
		Page:      page,
		PageSize:  pageSize,
	}

	result, err := mediatr.Send[*queries.GetAccountTransactionsQuery, *queries.GetAccountTransactionsResponse](c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ProcessTransaction godoc
// @Summary Process a transaction
// @Description Process a pending transaction to completion
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} commands.ProcessTransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions/{id}/process [post]
func (h *TransactionHandler) ProcessTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	cmd := &commands.ProcessTransactionCommand{ID: id}
	result, err := mediatr.Send[*commands.ProcessTransactionCommand, *commands.ProcessTransactionResponse](c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CancelTransaction godoc
// @Summary Cancel a transaction
// @Description Cancel a pending transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} commands.CancelTransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions/{id}/cancel [post]
func (h *TransactionHandler) CancelTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	cmd := &commands.CancelTransactionCommand{ID: id}
	result, err := mediatr.Send[*commands.CancelTransactionCommand, *commands.CancelTransactionResponse](c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
