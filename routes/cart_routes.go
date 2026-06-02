package routes

import (
	"demodiqit_api/controllers"
	"demodiqit_api/config"

	"github.com/gin-gonic/gin"
)

func SetupCartRoutes(api *gin.RouterGroup, _ *config.Config) {
	cartController := controllers.NewCartController()

	cart := api.Group("/cart")
	{
		cart.GET("", cartController.GetCart)
		cart.POST("", cartController.SaveCart)
	}
}
