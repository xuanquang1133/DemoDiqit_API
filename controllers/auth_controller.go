package controllers

import (
	"net/http"

	"demodiqit_api/config"
	"demodiqit_api/helpers/context"
	"demodiqit_api/helpers/crypt"
	"demodiqit_api/helpers/respond"
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
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid request payload",
			Code:    "AUTH-001",
		})
		return
	}

	var user models.User
	result := config.DB.Where("email = ? OR username = ?", req.Email, req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Message: "Invalid email/username or password",
			Code:    "AUTH-002",
		})
		return
	}

	// Compare password
	if !crypt.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Message: "Invalid email/username or password",
			Code:    "AUTH-002",
		})
		return
	}

	// Generate JWT Token
	token, err := crypt.GenerateJWT(user.ID, user.Username, user.Email, user.Roles, ac.Config.JWTSecret, ac.Config.JWTExpirationDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to generate token",
			Code:    "AUTH-003",
		})
		return
	}

	// Save token to user table
	user.UserToken = token
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to save user token",
			Code:    "AUTH-004",
		})
		return
	}

	rsp := respond.SuccessRespond{
		Message: "Login successfully!",
		Data: request.LoginResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Roles:    user.Roles,
			Token:    token,
		},
	}

	c.JSON(http.StatusOK, rsp)
}

// Me returns the profile of the currently authenticated user.
func (ac *AuthController) Me(c *gin.Context) {
	user := contextHelper.GetUserFromContext(c)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "AUTH-005",
			Message: "Unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "OK",
		Data: gin.H{
			"user_id":   user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"roles":     user.Roles,
		},
	})
}
