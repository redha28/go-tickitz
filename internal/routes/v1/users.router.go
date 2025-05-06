package v1

import (
	handler "gotickitz/internal/handlers"
	"gotickitz/internal/middlewares"
	"gotickitz/internal/repositories"

	"github.com/gin-gonic/gin"
)

func initUsersRouter(v1 *gin.RouterGroup, usersRepo *repositories.UsersRepository, middleWare *middlewares.Middleware) {
	users := v1.Group("/users")
	userHandler := handler.NewAuthHandler(usersRepo)
	profileHandler := handler.NewUsersHandler(usersRepo)
	// AUTH
	auth := users.Group("/auth")
	{
		auth.POST("", userHandler.LoginHandler)
		auth.POST("/new", userHandler.RegisterHandler)
		// auth.GET("/verify", userHandler.VerifyHandler)
	}

	// PROFILE
	profile := users.Group("/profile")
	{
		profile.GET("", middleWare.VerifyToken, middleWare.AccsessGate("user"), profileHandler.GetProfileHandler)
		profile.PATCH("", middleWare.VerifyToken, middleWare.AccsessGate("user"), profileHandler.UpdateProfileHandler)
	}

}
