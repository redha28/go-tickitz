package v1

import (
	handler "gotickitz/internal/handlers"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

func initMoviesRouter(v1 *gin.RouterGroup, moviesRepo *repositories.MoviesRepository) {
	// MOVIE
	movies := v1.Group("/movies")
	moviesHandler := handler.NewMoviesHandler(moviesRepo)
	{
		movies.GET("", moviesHandler.GetMoviesHandler)
		movies.GET("/:id", moviesHandler.GetMovieDetailHandler)
	}
}

// func getMoviesHandler(c *gin.Context)         {}
// func getMovieDetailHandler(c *gin.Context)    {}
