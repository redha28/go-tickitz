package v1

import "github.com/gin-gonic/gin"

func initAdminRouter(v1 *gin.RouterGroup) {
	// ADMIN
	admin := v1.Group("/admin")
	{
		admin.GET("/movies", getMoviesAdminHandler)
		admin.POST("/movies", createMovieHandler)
		admin.PATCH("/movies/:id", updateMovieHandler)
		admin.DELETE("/movies/:id", deleteMovieHandler)
	}
}

func getMoviesAdminHandler(c *gin.Context) {}
func createMovieHandler(c *gin.Context)    {}
func updateMovieHandler(c *gin.Context)    {}
func deleteMovieHandler(c *gin.Context)    {}
