package handler

import (
	"gotickitz/internal/models"
	"gotickitz/internal/repositories"
	"gotickitz/pkg"
	"log"
	"net/http"
	"regexp"
	"strings"

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
		if strings.Contains(err.Error(), "Field validation for 'Password'") {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Password length must be at least 8 characters",
			})
			return
		}
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid input",
		})
		return
	}
	if !isValidEmail(userReq.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid email format",
		})
		return
	}
	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	hashedPass, err := hash.GenHashedPassword(userReq.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Hash failed",
		})
		return
	}
	result, findEmail, err := a.UsersRepository.UseRegister(c.Request.Context(), userReq, hashedPass)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server error",
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
		"data": map[string]string{"email": result.Email, "role": result.Role},
	})
}

func loginHandler(c *gin.Context, a *AuthHandler) {
	var userReq models.UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		if strings.Contains(err.Error(), "Field validation for 'Password'") {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Password length must be at least 8 characters",
			})
			return
		}
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid input",
		})
		return
	}
	log.Println("[DEBUG] USERREQ:", userReq.Email)
	if !isValidEmail(userReq.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid email format",
		})
		return
	}
	result, err := a.UsersRepository.UseLogin(c.Request.Context(), userReq.Email)
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
	// log.Println("[DEBUG] PASSWORD:", result.Pass)
	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	valid, err := hash.CompareHashAndPassword(result.Pass, userReq.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Terjadi kesalahan pada server",
		})
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Email atau password salah",
		})
		return
	}
	payload := pkg.NewPayload(result.AuthID, result.Role)
	token, err := payload.GenerateToken()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server error",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "login success",
		"data": map[string]string{"id": result.AuthID, "role": result.Role, "token": token},
	})
}

func isValidEmail(email string) bool {
	// Menggunakan regex untuk memastikan format email benar
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
