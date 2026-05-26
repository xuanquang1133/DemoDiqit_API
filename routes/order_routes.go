package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"
	"demodiqit_api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	orderController := controllers.NewOrderController(cfg)

	orders := rg.Group("/orders")
	orders.Use(middleware.JWTAuthMiddleware(cfg))
	{
		orders.GET("", orderController.ListOrders)
		orders.GET("/:id", orderController.GetOrder)
		orders.POST("", orderController.CreateOrder)
		orders.PATCH("/:id/status", orderController.UpdateOrderStatus)
	}
}
