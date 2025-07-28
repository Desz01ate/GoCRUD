package router

import (
	"arise_tech_assetment/internal/api/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteRegistrar interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

func SetupRoutes(r *Router) {
	accountHandler := http.NewAccountHandler()
	transactionHandler := http.NewTransactionHandler()

	v1 := r.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			// HealthCheck godoc
			// @Summary Health check
			// @Description Get the health status of the service
			// @Tags health
			// @Accept json
			// @Produce json
			// @Success 200 {object} map[string]string
			// @Router /health [get]
			health.GET("", func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{
					"status":  "ok",
					"message": "Server is running",
				})
			})
		}

		accounts := v1.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("", accountHandler.GetAccounts)

			accounts.GET("/:id/transactions", transactionHandler.GetAccountTransactions)

			accounts.GET("/:id", accountHandler.GetAccount)
			accounts.PUT("/:id", accountHandler.UpdateAccount)
			accounts.DELETE("/:id", accountHandler.DeleteAccount)

			accounts.GET("/number/:number", accountHandler.GetAccountByNumber)
		}

		transactions := v1.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetTransactions)

			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.POST("/:id/process", transactionHandler.ProcessTransaction)
			transactions.POST("/:id/cancel", transactionHandler.CancelTransaction)
		}
	}

	r.Engine().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
