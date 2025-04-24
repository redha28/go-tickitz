package v1

import "github.com/gin-gonic/gin"

func initShowingsRouter(v1 *gin.RouterGroup) {
	// SHOWING MOVIE
	showings := v1.Group("/showings")
	{
		showings.GET("", getShowingsHandler)
		showings.GET("/:showing_id", getShowingDetailHandler)
	}
}

func getShowingsHandler(c *gin.Context)      {}
func getShowingDetailHandler(c *gin.Context) {}
