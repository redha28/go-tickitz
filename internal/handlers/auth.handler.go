package handler

import (
	"gotickitz/internal/models"
	"gotickitz/internal/repositories"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	*repositories.UsersRepository
}

func NewAuthHandler(authRepo *repositories.UsersRepository) *AuthHandler {
	return &AuthHandler{authRepo}
}

func (a *AuthHandler) LoginHandler(ctx *gin.Context) {
	loginHandler(ctx, a)
}

func (a *AuthHandler) RegisterHandler(ctx *gin.Context) {
	registerHandler(ctx, a)
}

func registerHandler(c *gin.Context, a *AuthHandler) {
	var userReq models.UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "Invalid input",
			"error": err.Error(),
		})
		return
	}
	result, findEmail, err := a.UsersRepository.UseRegister(c.Request.Context(), userReq)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Terjadi kesalahan pada server",
		})
		return
	}
	if findEmail.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Email sudah terdaftar",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"msg":  "success",
		"data": result,
	})
}

func loginHandler(c *gin.Context, a *AuthHandler) {
	var userReq models.UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "Invalid input",
			"error": err.Error(),
		})
		return
	}
	result, err := a.UsersRepository.UseLogin(c.Request.Context(), userReq.Email, userReq.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Terjadi kesalahan pada server",
		})
		return
	}
	if result.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Email atau password salah",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "login success",
		"data": result,
	})
}
