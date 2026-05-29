package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupPublicProductRoutes registers public product routes (no auth required)
func SetupPublicProductRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	productController := controllers.NewProductController(cfg)

	rg.GET("/products/slug/:slug", productController.GetProductBySlug)
}
