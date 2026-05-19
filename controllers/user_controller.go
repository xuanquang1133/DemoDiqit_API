package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"demodiqit_api/config"
	contextHelper "demodiqit_api/helpers/context"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Config *config.Config
}

func NewUserController(cfg *config.Config) *UserController {
	return &UserController{
		Config: cfg,
	}
}

// ListUser handles GET /users
func (uc *UserController) ListUser(c *gin.Context) {
	var req request.ListUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid query parameters",
			Code:    "USER-020",
		})
		return
	}

	var users []models.User
	query := config.DB.Model(&models.User{})

	if req.Keyword != "" {
		keyword := strings.ToLower(req.Keyword)
		query = query.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(full_name) LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if req.Role != "" && req.Role != "All" {
		query = query.Where("? = ANY(roles)", req.Role)
	}

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	var total int64
	query.Count(&total)

	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	if err := query.Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to fetch users",
			Code:    "USER-001",
		})
		return
	}

	userResponses := make([]request.UserResponse, 0)
	for _, u := range users {
		userResponses = append(userResponses, request.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FullName:  u.FullName,
			Roles:     u.Roles,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Success",
		Data: respond.PaginatedData{
			Items:      userResponses,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	})
}

// UserDetail handles GET /users/:id
func (uc *UserController) UserDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid user ID",
			Code:    "USER-002",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Message: "User not found",
			Code:    "USER-003",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Success",
		Data: request.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Roles:     user.Roles,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	})
}

// CreateUser handles POST /users
func (uc *UserController) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid request payload",
			Code:    "USER-004",
		})
		return
	}

	// Check if username or email exists
	var existingCount int64
	if err := config.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&existingCount).Error; err != nil || existingCount > 0 {
		c.JSON(http.StatusConflict, respond.ErrorRespond{
			Message: "Username or Email already exists",
			Code:    "USER-005",
		})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password, // Hook handles hashing
		Email:    req.Email,
		FullName: req.FullName,
		Roles:    req.Roles,
		IsActive: isActive,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to create user. Email or username might exist.",
			Code:    "USER-005",
		})
		return
	}

	c.JSON(http.StatusCreated, respond.SuccessRespond{
		Message: "User created successfully",
		Data: request.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Roles:     user.Roles,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	})
}

// UpdateUser handles PUT /users/:id
func (uc *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid user ID",
			Code:    "USER-006",
		})
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid request payload",
			Code:    "USER-007",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Message: "User not found",
			Code:    "USER-008",
		})
		return
	}

	currentUser := contextHelper.GetUserFromContext(c)
	if user.Username == "admin" && currentUser.Username != "admin" {
		c.JSON(http.StatusForbidden, respond.ErrorRespond{
			Message: "You are not allowed to modify the root admin account",
			Code:    "USER-017",
		})
		return
	}

	// Check if new username or email exists (excluding this user)
	var existingCount int64
	if err := config.DB.Model(&models.User{}).Where("(username = ? OR email = ?) AND id != ?", req.Username, req.Email, id).Count(&existingCount).Error; err != nil || existingCount > 0 {
		c.JSON(http.StatusConflict, respond.ErrorRespond{
			Message: "Username or Email already exists",
			Code:    "USER-EXISTS",
		})
		return
	}

	user.Username = req.Username
	user.Email = req.Email
	user.FullName = req.FullName
	user.Roles = req.Roles
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to update user",
			Code:    "USER-009",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "User updated successfully",
		Data: request.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Roles:     user.Roles,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	})
}

// DeleteUser handles DELETE /users/:id
func (uc *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid user ID",
			Code:    "USER-010",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Message: "User not found",
			Code:    "USER-011",
		})
		return
	}

	if user.Username == "admin" {
		c.JSON(http.StatusForbidden, respond.ErrorRespond{
			Message: "The root admin account cannot be deleted by anyone",
			Code:    "USER-018",
		})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to delete user",
			Code:    "USER-012",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "User deleted successfully",
		Data:    nil,
	})
}


// UpdateStatus handles PATCH /users/:id/status
func (uc *UserController) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid user ID",
			Code:    "USER-013",
		})
		return
	}

	var req request.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid request payload",
			Code:    "USER-014",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Message: "User not found",
			Code:    "USER-015",
		})
		return
	}

	if user.Username == "admin" {
		c.JSON(http.StatusForbidden, respond.ErrorRespond{
			Message: "The root admin account cannot be modified by anyone",
			Code:    "USER-019",
		})
		return
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to update user status",
			Code:    "USER-016",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "User status updated successfully",
		Data: request.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Roles:     user.Roles,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	})
}
