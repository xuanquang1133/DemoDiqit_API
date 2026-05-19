package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"
	"demodiqit_api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes registers all product-related routes
func SetupProductRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	productController := controllers.NewProductController(cfg)

	// Public routes (no authentication required)
	products := rg.Group("/products")
	{
		products.GET("", productController.ListProducts)
		products.GET("/:id", productController.GetProduct)
	}

	// Protected routes (admin only)
	protected := products.Group("")
	protected.Use(middleware.JWTAuthMiddleware(cfg))
	{
		protected.POST("", productController.CreateProduct)
		protected.PUT("/:id", productController.UpdateProduct)
		protected.DELETE("/:id", productController.DeleteProduct)
		protected.PATCH("/:id/status", productController.UpdateProductStatus)
	}
}
