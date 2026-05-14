package controllers

import (
	"net/http"

	"demodiqit_api/config"
	"demodiqit_api/helpers/crypt"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Config *config.Config
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		Config: cfg,
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload",
			"code":    "AUTH-001",
		})
		return
	}

	var user models.User
	result := config.DB.Where("email = ? OR username = ?", req.Email, req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email/username or password",
			"code":    "AUTH-002",
		})
		return
	}

	// Compare password
	if !crypt.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email/username or password",
			"code":    "AUTH-002",
		})
		return
	}

	// Generate JWT Token
	token, err := crypt.GenerateJWT(user.ID, user.Username, user.Email, ac.Config.JWTSecret, ac.Config.JWTExpirationDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
			"code":    "AUTH-003",
		})
		return
	}

	// Save token to user table
	user.UserToken = token
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save user token",
			"code":    "AUTH-004",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"roles":     user.Roles,
		},
	})
}
