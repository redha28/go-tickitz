package v1

import (
	handler "gotickitz/internal/handlers"
	"gotickitz/internal/middlewares"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

func initAdminRouter(v1 *gin.RouterGroup, adminRepo *repositories.AdminRepository, middleWare *middlewares.Middleware) {
	// ADMIN
	admin := v1.Group("/admin")
	adminHandler := handler.NewAdminHandler(adminRepo)
	// MOVIE
	{
		// admin.GET("/movies", getMoviesAdminHandler)
		admin.POST("/movies", middleWare.VerifyToken, middleWare.AccsessGate("admin"), adminHandler.CreateMovieHandler)
		admin.PATCH("/movies/:id", middleWare.VerifyToken, middleWare.AccsessGate("admin"), adminHandler.UpdateMovieHandler)
		admin.DELETE("/movies", middleWare.VerifyToken, middleWare.AccsessGate("admin"), adminHandler.DeleteMovieHandler)
	}
}

// func getMoviesAdminHandler(c *gin.Context) {}

// func createMovieHandler(c *gin.Context)    {}
// func updateMovieHandler(c *gin.Context) {}
// func deleteMovieHandler(c *gin.Context) {}
