package routes

import (
	v1 "gotickitz/internal/routes/v1"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(pg *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	v1.InitRouter(router, pg)
	return router
}
