package v1

import (
	handler "gotickitz/internal/handlers"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

func initShowingsRouter(v1 *gin.RouterGroup, showingRepo *repositories.ShowingsRepository) {
	// SHOWING MOVIE
	showings := v1.Group("/showings")
	showingsHandler := handler.NewShowingsHandler(showingRepo)
	{
		showings.GET("", showingsHandler.GetShowingsHandler)
		showings.GET("/:id/seat", showingsHandler.GetSeatsHandler)
	}
}

// func getShowingsHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"msg": "Bentar yaa Cinema tutup"})
// }

// func getShowingDetailHandler(c *gin.Context) {}
