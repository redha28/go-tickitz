package handler

import (
	"gotickitz/internal/models"
	"gotickitz/internal/repositories"
	"gotickitz/pkg"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	*repositories.AdminRepository
}

func NewAdminHandler(adminRepo *repositories.AdminRepository) *AdminHandler {
	return &AdminHandler{adminRepo}
}

func (a *AdminHandler) CreateMovieHandler(ctx *gin.Context) {
	var movieReq models.AdminCreateMovieReq
	payloads, _ := ctx.Get("payloads")
	userPayload := payloads.(*pkg.Payload)
	if err := ctx.ShouldBindJSON(&movieReq); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	_, err := time.Parse("2006-01-02", movieReq.Release)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	// isAdmin, err := a.UseCheckAdmin(ctx.Request.Context(), userPayload.Id)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
	// 	return
	// }
	// if !isAdmin {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Only admin can create movie"})
	// 	return
	// }
	err = a.UseCreateMovie(ctx.Request.Context(), movieReq, userPayload.Id)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (a *AdminHandler) UpdateMovieHandler(ctx *gin.Context) {
	idMovie, ok := ctx.Params.Get("id")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	var movieReq models.AdminEditMovieReq
	if err := ctx.ShouldBindJSON(&movieReq); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	idMovieInt, err := strconv.Atoi(idMovie)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	msg, err := a.UseEditMovie(ctx.Request.Context(), movieReq, idMovieInt)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	if msg == "no fields to update" && len(movieReq.Genres) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Require at least one field to update"})
		return
	}
	if len(movieReq.Genres) != 0 {
		errorGenres := a.UseEditMovieGenre(ctx.Request.Context(), idMovieInt, movieReq.Genres)
		if errorGenres != nil {
			log.Println(errorGenres.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (a *AdminHandler) DeleteMovieHandler(ctx *gin.Context) {
	var movieReq models.AdminDeleteMovieReq
	if err := ctx.ShouldBindJSON(&movieReq); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	err := a.UseDeleteMovie(ctx.Request.Context(), movieReq)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
