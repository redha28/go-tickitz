package handler

import (
	"log"
	"net/http"
	"strconv"

	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

type MoviesHandler struct {
	moviesRepo *repositories.MoviesRepository
}

func NewMoviesHandler(moviesRepo *repositories.MoviesRepository) *MoviesHandler {
	return &MoviesHandler{moviesRepo}
}

func (m *MoviesHandler) GetMoviesHandler(ctx *gin.Context) {

	name := ctx.Query("name")
	limitStr := ctx.Query("limit")
	offsetStr := ctx.Query("offset")
	options := ctx.Query("options")
	genre := ctx.Query("genre")

	limit := 5
	offset := 0

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid limit value"})
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid offset value"})
			return
		}
		offset = parsedOffset
	}

	if options == "" {
		options = "default"
	}

	if options != "allrelease" && options != "upcoming" && options != "popular" && options != "default" {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid options value"})
		return
	}

	movies, meta, err := m.moviesRepo.UseGetMovies(ctx, name, options, genre, limit, offset)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Failed to fetch movies",
		})
		return
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "No movies found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": movies,
		"meta": meta,
	})
}

func (m *MoviesHandler) GetMovieDetailHandler(ctx *gin.Context) {
	idReq := ctx.Param("id")

	if idReq == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}
	id, err := strconv.Atoi(idReq)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}
	movie, err := m.moviesRepo.UseGetMovie(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch movie",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": movie,
	})
}
