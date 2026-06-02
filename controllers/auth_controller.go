package controllers

import (
	"net/http"

	"demodiqit_api/config"
	contextHelper "demodiqit_api/helpers/context"
	"demodiqit_api/helpers/crypt"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type AuthController struct {
	Config *config.Config
}

// Register handles POST /auth/register — creates a new customer account
func (ac *AuthController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "AUTH-010",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Check duplicate
	var count int64
	config.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, respond.ErrorRespond{
			Code:    "AUTH-011",
			Message: "Username or email already exists",
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password, // BeforeCreate hook hashes it
		Email:    req.Email,
		FullName: req.FullName,
		Roles:    pq.StringArray{"customer"},
		IsActive: true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "AUTH-012",
			Message: "Failed to create account",
		})
		return
	}

	c.JSON(http.StatusCreated, respond.SuccessRespond{
		Message: "Account created successfully",
		Data: map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"full_name": user.FullName,
		},
	})
}

// ChangePassword handles PUT /auth/change-password — changes the current user's password
func (ac *AuthController) ChangePassword(c *gin.Context) {
	user := contextHelper.GetUserFromContext(c)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "AUTH-005",
			Message: "Unauthorized",
		})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword    string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "AUTH-013",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	var dbUser models.User
	if err := config.DB.First(&dbUser, user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "AUTH-014",
			Message: "User not found",
		})
		return
	}

	if !crypt.CheckPasswordHash(req.CurrentPassword, dbUser.Password) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "AUTH-015",
			Message: "Current password is incorrect",
		})
		return
	}

	hashed, err := crypt.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "AUTH-016",
			Message: "Failed to hash new password",
		})
		return
	}

	dbUser.Password = hashed
	if err := config.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "AUTH-017",
			Message: "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Password changed successfully",
	})
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
	result := config.DB.Where("is_active = ? AND (email = ? OR username = ?)", true, req.Email, req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid email/username or password",
			Code:    "AUTH-002",
		})
		return
	}

	// Compare password
	if !crypt.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
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

// UserInfoByToken returns the profile of the currently authenticated user.
func (ac *AuthController) UserInfoByToken(c *gin.Context) {
	user := contextHelper.GetUserFromContext(c)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "AUTH-005",
			Message: "Unauthorized",
		})
		return
	}

	rsp := respond.SuccessRespond{
		Message: "OK",
		Data: request.UserInfoByTokenResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Roles:    user.Roles,
		},
	}

	c.JSON(http.StatusOK, rsp)
}
