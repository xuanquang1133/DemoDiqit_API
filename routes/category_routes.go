package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	categoryController := controllers.NewCategoryController(cfg)

	categories := rg.Group("/categories")
	{
		categories.GET("", categoryController.ListCategory)
		categories.GET("/public", categoryController.ListCommon)
		categories.GET("/:id", categoryController.CategoryDetail)
		categories.POST("", categoryController.CreateCategory)
		categories.PUT("/:id", categoryController.UpdateCategory)
		categories.DELETE("/:id", categoryController.DeleteCategory)
		categories.PATCH("/:id/status", categoryController.UpdateStatus)
	}
}
