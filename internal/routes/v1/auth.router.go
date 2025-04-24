package v1

import (
	handler "gotickitz/internal/handlers"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

func initAuthRouter(v1 *gin.RouterGroup, usersRepo *repositories.UsersRepository) {
	users := v1.Group("/users")
	userHandler := handler.NewAuthHandler(usersRepo)
	profileHandler := handler.NewUsersHandler(usersRepo)
	// AUTH
	auth := users.Group("/auth")
	{
		auth.GET("", userHandler.LoginHandler)
		auth.POST("", userHandler.RegisterHandler)
	}

	// PROFILE
	profile := users.Group("/profile")
	{
		profile.GET("/:uuid", profileHandler.GetProfileHandler)
		profile.PATCH("/:uuid", profileHandler.UpdateProfileHandler)
	}

}
