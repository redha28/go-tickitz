package v1

import (
	handler "gotickitz/internal/handlers"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

func initTransactionsRouter(v1 *gin.RouterGroup, transactionRepo *repositories.TransactionsRepository) {

	// TRANSACTIONS
	transactions := v1.Group("/transactions")
	transactionsHandler := handler.NewTransactionsHandler(transactionRepo)
	{
		transactions.GET("/:user_id", getTransactionsHandler)
		transactions.POST("", transactionsHandler.CreateTransactionHandler)
	}
}

func getTransactionsHandler(c *gin.Context) {}

// func createTransactionHandler(c *gin.Context) {}
