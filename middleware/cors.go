package middleware

import (
	"demodiqit_api/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsConfig returns the CORS middleware for Gin
func CorsConfig(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
