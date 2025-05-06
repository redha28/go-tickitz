package v1

import (
	"gotickitz/internal/middlewares"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitRouter(router *gin.Engine, pg *pgxpool.Pool, rdb *redis.Client) {
	v1 := router.Group("/api/v1")

	usersRepo := repositories.NewUsersRepository(pg)
	moviesRepo := repositories.NewMoviesRepository(pg, rdb)
	showingsRepo := repositories.NewShowingsRepository(pg)
	transactionsRepo := repositories.NewTransactionsRepository(pg)
	adminRepo := repositories.NewAdminRepository(pg, rdb)

	// Middleware
	middlewares := middlewares.InitMiddleware()

	initAdminRouter(v1, adminRepo, middlewares)
	initUsersRouter(v1, usersRepo, middlewares)
	initMoviesRouter(v1, moviesRepo)
	initShowingsRouter(v1, showingsRepo)
	initTransactionsRouter(v1, transactionsRepo, middlewares)
}
