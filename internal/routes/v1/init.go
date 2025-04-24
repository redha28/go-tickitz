package v1

import (
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(router *gin.Engine, pg *pgxpool.Pool) {
	v1 := router.Group("/api/v1")

	usersRepo := repositories.NewUsersRepository(pg)

	initAdminRouter(v1)
	initAuthRouter(v1, usersRepo)
	initMoviesRouter(v1)
	initShowingsRouter(v1)
	initTransactionsRouter(v1)
}
