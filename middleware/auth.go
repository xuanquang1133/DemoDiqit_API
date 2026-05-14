package middleware

import (
	"net/http"
	"strings"

	"demodiqit_api/config"
	"demodiqit_api/helpers/context"
	"demodiqit_api/helpers/crypt"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/models"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates the Bearer token in the Authorization header.
// On success, it fetches the user from DB and stores UserContext in the Gin context.
func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, respond.ErrorRespond{
				Code:    "AUTH-005",
				Message: "Authorization header is missing or malformed",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := crypt.ParseJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, respond.ErrorRespond{
				Code:    "AUTH-006",
				Message: "Invalid or expired token",
			})
			return
		}

		var user models.User
		if err := config.DB.Where("user_token = ?", tokenString).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, respond.ErrorRespond{
				Code:    "AUTH-007",
				Message: "User not found or token revoked",
			})
			return
		}

		userCtx := contextHelper.UserContext{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Roles:    user.Roles,
		}

		// Store user context for use by subsequent handlers/middleware
		c.Set(contextHelper.UserContextKey, userCtx)
		c.Next()
	}
}
