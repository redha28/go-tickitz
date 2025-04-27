package handler

import (
	"gotickitz/internal/models"
	"gotickitz/internal/repositories"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	*repositories.UsersRepository
}

func NewUsersHandler(usersRepo *repositories.UsersRepository) *UsersHandler {
	return &UsersHandler{usersRepo}
}

func (u *UsersHandler) GetProfileHandler(c *gin.Context) {
	var req models.IdParams
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query := `SELECT p.firstname, p.lastname, p.picture, p.phone, p.point, a.email
	          FROM profile p
	          JOIN auth a ON p.auth_id = a.id
	          WHERE p.auth_id = $1`
	var profile models.ProfileRes
	if err := u.QueryRow(c.Request.Context(), query, req.UUID).Scan(&profile.Firstname, &profile.Lastname, &profile.Picture, &profile.Phone, &profile.Point, &profile.Email); err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "Profile not found",
			})
			return
		}
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "Profile retrieved successfully",
		"data": profile,
	})
}

func (u *UsersHandler) UpdateProfileHandler(c *gin.Context) {
	var req models.IdParams
	if err := c.ShouldBindUri(&req); err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{"msg": "Invalid input"})
		return
	}
	var updateReq models.UpdateProfileReq
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid input",
		})
		return
	}
	err := u.UseUpdateProfile(c.Request.Context(), req.UUID, updateReq)
	if err != nil {
		if err.Error() == "no fields to update" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "No fields to update",
			})
			return
		}
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Profile updated successfully",
	})
}
