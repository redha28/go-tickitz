package v1

import (
	handler "gotickitz/internal/handlers"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

func initAdminRouter(v1 *gin.RouterGroup, adminRepo *repositories.AdminRepository) {
	// ADMIN
	admin := v1.Group("/admin")
	adminHandler := handler.NewAdminHandler(adminRepo)
	// MOVIE
	{
		admin.GET("/movies", getMoviesAdminHandler)
		admin.POST("/movies", adminHandler.CreateMovieHandler)
		admin.PATCH("/movies/:id", adminHandler.UpdateMovieHandler)
		admin.DELETE("/movies", adminHandler.DeleteMovieHandler)
	}
}

func getMoviesAdminHandler(c *gin.Context) {}

// func createMovieHandler(c *gin.Context)    {}
// func updateMovieHandler(c *gin.Context) {}
// func deleteMovieHandler(c *gin.Context) {}
