package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"demodiqit_api/config"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	Config *config.Config
}

func NewCategoryController(cfg *config.Config) *CategoryController {
	return &CategoryController{
		Config: cfg,
	}
}

// ListCategory handles GET /categories
func (cc *CategoryController) ListCategory(c *gin.Context) {
	var req request.ListCategoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid query parameters",
			Code:    "CATEGORY-020",
		})
		return
	}

	var categories []models.Category
	query := config.DB.Model(&models.Category{})

	if req.Keyword != "" {
		keyword := strings.ToLower(req.Keyword)
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%")
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

	if err := query.Offset(offset).Limit(limit).Order("id desc").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to fetch categories",
			Code:    "CATEGORY-001",
		})
		return
	}

	categoryResponses := make([]request.CategoryResponse, 0)
	for _, cat := range categories {
		categoryResponses = append(categoryResponses, request.CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			Code:      cat.Code,
			IsActive:  cat.IsActive,
			CreatedAt: cat.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Success",
		Data: respond.PaginatedData{
			Items:      categoryResponses,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	})
}

// CategoryDetail handles GET /categories/:id
func (cc *CategoryController) CategoryDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid category ID",
			Code:    "CATEGORY-002",
		})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Message: "Category not found",
			Code:    "CATEGORY-003",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Success",
		Data: request.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Code:      category.Code,
			IsActive:  category.IsActive,
			CreatedAt: category.CreatedAt,
		},
	})
}

// CreateCategory handles POST /categories
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid request payload",
			Code:    "CATEGORY-004",
		})
		return
	}

	// Check if name or code exists (excluding soft-deleted)
	var existingCount int64
	if err := config.DB.Model(&models.Category{}).Where("(name = ? OR code = ?) AND deleted_at IS NULL", req.Name, req.Code).Count(&existingCount).Error; err != nil || existingCount > 0 {
		c.JSON(http.StatusConflict, respond.ErrorRespond{
			Message: "Name or Code already exists",
			Code:    "CATEGORY-005",
		})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category := models.Category{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: isActive,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to create category",
			Code:    "CATEGORY-005",
		})
		return
	}

	c.JSON(http.StatusCreated, respond.SuccessRespond{
		Message: "Category created successfully",
		Data: request.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Code:      category.Code,
			IsActive:  category.IsActive,
			CreatedAt: category.CreatedAt,
		},
	})
}

// UpdateCategory handles PUT /categories/:id
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid category ID",
			Code:    "CATEGORY-006",
		})
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid request payload",
			Code:    "CATEGORY-007",
		})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Message: "Category not found",
			Code:    "CATEGORY-008",
		})
		return
	}

	// Check if new name or code exists (excluding soft-deleted and this category)
	var existingCount int64
	if err := config.DB.Model(&models.Category{}).Where("(name = ? OR code = ?) AND id != ? AND deleted_at IS NULL", req.Name, req.Code, id).Count(&existingCount).Error; err != nil || existingCount > 0 {
		c.JSON(http.StatusConflict, respond.ErrorRespond{
			Message: "Name or Code already exists",
			Code:    "CATEGORY-EXISTS",
		})
		return
	}

	category.Name = req.Name
	category.Code = req.Code
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to update category",
			Code:    "CATEGORY-009",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Category updated successfully",
		Data: request.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Code:      category.Code,
			IsActive:  category.IsActive,
			CreatedAt: category.CreatedAt,
		},
	})
}

// DeleteCategory handles DELETE /categories/:id
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid category ID",
			Code:    "CATEGORY-010",
		})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Message: "Category not found",
			Code:    "CATEGORY-011",
		})
		return
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to delete category",
			Code:    "CATEGORY-012",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Category deleted successfully",
		Data:    nil,
	})
}


// UpdateStatus handles PATCH /categories/:id/status
func (cc *CategoryController) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid category ID",
			Code:    "CATEGORY-013",
		})
		return
	}

	var req request.UpdateCategoryStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid request payload",
			Code:    "CATEGORY-014",
		})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Message: "Category not found",
			Code:    "CATEGORY-015",
		})
		return
	}

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to update category status",
			Code:    "CATEGORY-016",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Category status updated successfully",
		Data: request.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Code:      category.Code,
			IsActive:  category.IsActive,
			CreatedAt: category.CreatedAt,
		},
	})
}

// ListCommon handles GET /categories/list-common
func (cc *CategoryController) ListCommon(c *gin.Context) {
	var req request.ListCommonCategoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid query parameters",
			Code:    "CATEGORY-018",
		})
		return
	}

	var categories []models.Category
	query := config.DB.Model(&models.Category{})

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if err := query.Order("name asc").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to fetch categories",
			Code:    "CATEGORY-017",
		})
		return
	}

	categoryResponses := make([]request.CategoryResponse, 0)
	for _, cat := range categories {
		categoryResponses = append(categoryResponses, request.CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			Code:      cat.Code,
			IsActive:  cat.IsActive,
			CreatedAt: cat.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Success",
		Data:    categoryResponses,
	})
}
