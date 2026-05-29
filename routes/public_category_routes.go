package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupPublicCategoryRoutes registers public category routes (no auth required)
func SetupPublicCategoryRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	categoryController := controllers.NewCategoryController(cfg)

	rg.GET("/categories/list-common", categoryController.ListCommon)
}
