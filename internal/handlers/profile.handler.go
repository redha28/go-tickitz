package handler

import (
	"gotickitz/internal/models"
	"gotickitz/internal/repositories"
	"gotickitz/internal/utils"
	"gotickitz/pkg"
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
	payloads, _ := c.Get("payloads")
	userPayload := payloads.(*pkg.Payload)
	profile, err := u.UseGetProfile(c.Request.Context(), userPayload.Id)
	if err != nil {
		if err.Error() == "profile not found" {
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
	payloads, _ := c.Get("payloads")
	userPayload := payloads.(*pkg.Payload)
	var updateReq models.UpdateProfileReq
	if err := c.ShouldBind(&updateReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid input",
		})
		return
	}
	file := updateReq.Picture
	var filename string
	if file != nil {
		var err error
		oldFilename, err := u.UseGetProfile(c.Request.Context(), userPayload.Id)
		if err != nil {
			if err.Error() == "profile not found" {
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
		utils := utils.InitUtils()
		filename, _, err = utils.FileHandling(c, file, userPayload, oldFilename.Picture)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Terjadi kesalahan upload",
			})
		}
	}
	err := u.UseUpdateProfile(c.Request.Context(), userPayload.Id, updateReq, filename)
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
