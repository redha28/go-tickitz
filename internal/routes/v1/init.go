package v1

import (
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(router *gin.Engine, pg *pgxpool.Pool) {
	v1 := router.Group("/api/v1")

	usersRepo := repositories.NewUsersRepository(pg)
	moviesRepo := repositories.NewMoviesRepository(pg)
	showingsRepo := repositories.NewShowingsRepository(pg)
	transactionsRepo := repositories.NewTransactionsRepository(pg)
	adminRepo := repositories.NewAdminRepository(pg)

	initAdminRouter(v1, adminRepo)
	initAuthRouter(v1, usersRepo)
	initMoviesRouter(v1, moviesRepo)
	initShowingsRouter(v1, showingsRepo)
	initTransactionsRouter(v1, transactionsRepo)
}
