package v1

import "github.com/gin-gonic/gin"

func initTransactionsRouter(v1 *gin.RouterGroup) {

	// TRANSACTIONS
	transactions := v1.Group("/transactions")
	{
		transactions.GET("/:user_id", getTransactionsHandler)
		transactions.POST("", createTransactionHandler)
	}
}

func getTransactionsHandler(c *gin.Context)   {}
func createTransactionHandler(c *gin.Context) {}
