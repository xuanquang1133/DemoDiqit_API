package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	userController := controllers.NewUserController(cfg)

	users := rg.Group("/users")
	{
		users.GET("", userController.ListUser)
		users.GET("/:id", userController.UserDetail)
		users.POST("", userController.CreateUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
		users.PATCH("/:id/status", userController.UpdateStatus)
	}
}
