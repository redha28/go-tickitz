package v1

import "github.com/gin-gonic/gin"

func initMoviesRouter(v1 *gin.RouterGroup) {
	// MOVIE
	movies := v1.Group("/movies")
	{
		movies.GET("", getMoviesHandler)
		movies.GET("/:id", getMovieDetailHandler)
		movies.GET("/popular", getPopularMoviesHandler)
		movies.GET("/upcoming", getUpcomingMoviesHandler)
	}
}

func getMoviesHandler(c *gin.Context)         {}
func getMovieDetailHandler(c *gin.Context)    {}
func getPopularMoviesHandler(c *gin.Context)  {}
func getUpcomingMoviesHandler(c *gin.Context) {}
