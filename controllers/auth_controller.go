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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	var user models.User
	result := config.DB.Where("email = ? OR username = ?", req.Email, req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Tài khoản hoặc mật khẩu không chính xác!"})
		return
	}

	// Compare password
	if !crypt.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Tài khoản hoặc mật khẩu không chính xác!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"roles":     user.Roles,
		},
	})
}
