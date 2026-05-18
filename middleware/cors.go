package middleware

import (
	"demodiqit_api/config"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsConfig returns the CORS middleware for Gin
func CorsConfig(cfg *config.Config) gin.HandlerFunc {
	var feOrigins []string
	strs := strings.Split(cfg.FrontendOrigin, ";")
	for _, feOrigin := range strs {
		feOrigins = append(feOrigins, feOrigin)
	}

	return cors.New(cors.Config{
		AllowOrigins:     feOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
