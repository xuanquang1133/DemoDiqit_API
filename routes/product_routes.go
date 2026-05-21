package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"
	"demodiqit_api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	productController := controllers.NewProductController(cfg)

	products := rg.Group("/products")
	products.Use(middleware.JWTAuthMiddleware(cfg))
	{
		products.GET("", productController.ListProducts)
		products.GET("/:id", productController.GetProduct)
		products.POST("", productController.CreateProduct)
		products.PUT("/:id", productController.UpdateProduct)
		products.DELETE("/:id", productController.DeleteProduct)
		products.PATCH("/:id/status", productController.UpdateProductStatus)
	}
}
