package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	authController := controllers.NewAuthController(cfg)

	auth := rg.Group("/auth")
	{
		auth.POST("/login", authController.Login)
	}
}
