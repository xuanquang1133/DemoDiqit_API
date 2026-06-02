package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupGuestOrderRoutes(r *gin.Engine, cfg *config.Config) {
	guestGroup := r.Group("/api/v1/guest")
	{
		guestGroup.POST("/orders", controllers.GuestCreateOrder)
	}
}
