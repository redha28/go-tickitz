package handler

import (
	"gotickitz/internal/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShowingsHandler struct {
	*repositories.ShowingsRepository
}

func NewShowingsHandler(showingsRepo *repositories.ShowingsRepository) *ShowingsHandler {
	return &ShowingsHandler{showingsRepo}
}

func (s *ShowingsHandler) GetShowingsHandler(c *gin.Context) {
	id := c.Query("movie_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid movie ID"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid movie ID"})
		return
	}

	showings, err := s.UseGetShowings(c.Request.Context(), idInt)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch showings"})
		return
	}

	if len(showings) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Showings not found"})
		return
	}

	formattedShowings := make([]map[string]any, 0)
	for _, showing := range showings {
		formattedShowings = append(formattedShowings, map[string]any{
			"id":          showing.ID,
			"date":        showing.Date.Format("2006-01-02"),
			"time":        showing.Time.Format("15:04:05"),
			"cinemaName":  showing.CinemaName,
			"cinemaImage": showing.CinemaImage,
			"cityName":    showing.CityName,
			"price":       showing.Price,
			"movieTitle":  showing.MovieTitle,
		})
	}

	c.JSON(http.StatusOK, formattedShowings)
}

func (s *ShowingsHandler) GetSeatsHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid showing ID"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid showing ID"})
		return
	}

	seats, err := s.UseGetSeats(c.Request.Context(), idInt)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch seats"})
		return
	}

	if len(seats) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Seats not found"})
		return
	}

	c.JSON(http.StatusOK, seats)
}
