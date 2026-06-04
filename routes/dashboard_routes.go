package routes

import (
	"demodiqit_api/config"
	"demodiqit_api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupDashboardRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	rg.GET("/dashboard", controllers.GetDashboardStats)
	rg.GET("/dashboard/v2/chart", controllers.GetDashboardChart)
}
