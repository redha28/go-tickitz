package routes

import (
	"gotickitz/internal/middlewares"
	v1 "gotickitz/internal/routes/v1"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitRouter(pg *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.InitMiddleware().CORSMiddleware)
	router.Static("/public", "./public")
	v1.InitRouter(router, pg, rdb)
	return router
}
