package contextHelper

import (
	"github.com/gin-gonic/gin"
)

// UserContext holds necessary user information for the request context
type UserContext struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	FullName string   `json:"full_name"`
	Roles    []string `json:"roles"`
}

const UserContextKey = "userContext"

// GetUserFromContext retrieves the UserContext from the Gin context.
// Returns an empty UserContext if not found or invalid type.
func GetUserFromContext(c *gin.Context) UserContext {
	val, exists := c.Get(UserContextKey)
	if !exists {
		return UserContext{}
	}

	userCtx, ok := val.(UserContext)
	if !ok {
		return UserContext{}
	}

	return userCtx
}
