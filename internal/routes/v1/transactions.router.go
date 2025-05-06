package v1

import (
	handler "gotickitz/internal/handlers"
	"gotickitz/internal/middlewares"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

func initTransactionsRouter(v1 *gin.RouterGroup, transactionRepo *repositories.TransactionsRepository, middleWare *middlewares.Middleware) {

	// TRANSACTIONS
	transactions := v1.Group("/transactions")
	transactionsHandler := handler.NewTransactionsHandler(transactionRepo)
	{
		transactions.GET("", middleWare.VerifyToken, middleWare.AccsessGate("user"), transactionsHandler.GetTransactionsHandler)
		transactions.POST("", middleWare.VerifyToken, middleWare.AccsessGate("user"), transactionsHandler.CreateTransactionHandler)
	}
}

// func getTransactionsHandler(c *gin.Context) {}

// func createTransactionHandler(c *gin.Context) {}
