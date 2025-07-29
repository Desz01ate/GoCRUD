package http

import (
	"arise_tech_assessment/internal/application/commands"
	"arise_tech_assessment/internal/application/queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

type AccountHandler struct {
}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{}
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account with holder name, number, and initial balance
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body commands.CreateAccountCommand true "Account creation data"
// @Success 201 {object} commands.CreateAccountResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var cmd commands.CreateAccountCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := mediatr.Send[*commands.CreateAccountCommand, *commands.CreateAccountResponse](c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetAccount godoc
// @Summary Get account by ID
// @Description Get a single account by its ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} queries.GetAccountResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	query := &queries.GetAccountQuery{ID: id}
	result, err := mediatr.Send[*queries.GetAccountQuery, *queries.GetAccountResponse](c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAccounts godoc
// @Summary Get all accounts
// @Description Get a paginated list of all accounts
// @Tags accounts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} queries.GetAccountsResponse
// @Failure 500 {object} map[string]string
// @Router /accounts [get]
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := &queries.GetAccountsQuery{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := mediatr.Send[*queries.GetAccountsQuery, *queries.GetAccountsResponse](c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAccountByNumber godoc
// @Summary Get account by number
// @Description Get a single account by its account number
// @Tags accounts
// @Accept json
// @Produce json
// @Param number path string true "Account Number"
// @Success 200 {object} queries.GetAccountByNumberResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/number/{number} [get]
func (h *AccountHandler) GetAccountByNumber(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account number is required"})
		return
	}

	query := &queries.GetAccountByNumberQuery{Number: number}
	result, err := mediatr.Send[*queries.GetAccountByNumberQuery, *queries.GetAccountByNumberResponse](c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateAccount godoc
// @Summary Update an account
// @Description Update an existing account's information
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param account body commands.UpdateAccountCommand true "Account update data"
// @Success 200 {object} commands.UpdateAccountResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var cmd commands.UpdateAccountCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ID = id
	result, err := mediatr.Send[*commands.UpdateAccountCommand, *commands.UpdateAccountResponse](c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteAccount godoc
// @Summary Delete an account
// @Description Delete an existing account by ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} commands.DeleteAccountResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	cmd := &commands.DeleteAccountCommand{ID: id}
	result, err := mediatr.Send[*commands.DeleteAccountCommand, *commands.DeleteAccountResponse](c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
