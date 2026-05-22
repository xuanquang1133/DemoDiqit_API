package controllers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"demodiqit_api/config"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/helpers/utils"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
)

// ProductController handles product-related HTTP requests
type ProductController struct {
	cfg *config.Config
}

// NewProductController creates a new ProductController instance
func NewProductController(cfg *config.Config) *ProductController {
	return &ProductController{cfg: cfg}
}

// ListProducts handles GET /products
// Returns a paginated list of products with optional keyword filtering
func (pc *ProductController) ListProducts(c *gin.Context) {
	var query request.ProductListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-001",
			Message: "Invalid query parameters",
		})
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	offset := (query.Page - 1) * query.Limit

	db := config.DB.Model(&models.Product{}).Preload("Category")

	if query.Keyword != "" {
		keyword := "%" + strings.ToLower(query.Keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ?", keyword, keyword)
	}

	if query.IsCategory != "" {
		catID, err := strconv.ParseUint(query.IsCategory, 10, 32)
		if err == nil {
			db = db.Where("category_id = ?", uint(catID))
		}
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "PRD-002",
			Message: "Failed to count products",
		})
		return
	}

	var products []models.Product
	if err := db.Order("created_at DESC").Offset(offset).Limit(query.Limit).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "PRD-003",
			Message: "Failed to fetch products",
		})
		return
	}

	items := make([]request.ProductResponse, len(products))
	for i, p := range products {
		items[i] = toProductResponse(p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))

	c.JSON(http.StatusOK, respond.PaginatedData{
		Items:      items,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
	})
}

// GetProduct handles GET /products/:id
// Returns details of a single product by ID
func (pc *ProductController) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-004",
			Message: "Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := config.DB.Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "PRD-005",
			Message: "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(product))
}

// CreateProduct handles POST /products
// Creates a new product
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-006",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-007",
			Message: "Invalid price value",
		})
		return
	}

	var product models.Product
	cleanSlug := utils.Slug(req.Slug)
	if cleanSlug != "" && product.ExistsBySlug(config.DB, cleanSlug) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-026",
			Message: "Slug already exists. Please use a different product name or enter a custom slug.",
		})
		return
	}

	cleanSKU := utils.SKU(req.SKU)
	if cleanSKU != "" && product.ExistsBySKU(config.DB, cleanSKU) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-020",
			Message: "SKU already exists",
		})
		return
	}

	salePrice := float64(0)
	if req.SalePrice != nil {
		salePrice = *req.SalePrice
	}

	product = models.Product{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Slug:       cleanSlug,
		SKU:        cleanSKU,
		Description: req.Description,
		Price:      req.Price,
		SalePrice:  salePrice,
		Thumbnail:  req.Thumbnail,
		IsActive:   true,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "PRD-009",
			Message: "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(product))
}

// UpdateProduct handles PUT /products/:id
// Updates an existing product
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-010",
			Message: "Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "PRD-011",
			Message: "Product not found",
		})
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-012",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}

	if req.Slug != "" {
		cleanSlug := utils.Slug(req.Slug)
		if product.ExistsBySlug(config.DB, cleanSlug, product.ID) {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "PRD-018",
				Message: "Slug already exists",
			})
			return
		}
		product.Slug = cleanSlug
	}

	if req.SKU != "" {
		cleanSKU := utils.SKU(req.SKU)
		if product.ExistsBySKU(config.DB, cleanSKU, product.ID) {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "PRD-021",
				Message: "SKU already exists",
			})
			return
		}
		product.SKU = cleanSKU
	}

	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.SalePrice != nil && *req.SalePrice > 0 {
		product.SalePrice = *req.SalePrice
	}
	if req.Thumbnail != "" {
		product.Thumbnail = req.Thumbnail
	}

	// CategoryID: pointer — allow unsetting (nil) or setting
	product.CategoryID = req.CategoryID

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "PRD-015",
			Message: "Failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(product))
}

// DeleteProduct handles DELETE /products/:id
// Soft deletes a product
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-016",
			Message: "Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "PRD-017",
			Message: "Product not found",
		})
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "PRD-019",
			Message: "Failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Product deleted successfully",
	})
}

// UpdateProductStatus handles PATCH /products/:id/status
// Updates only the is_active field of a product
func (pc *ProductController) UpdateProductStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-022",
			Message: "Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "PRD-023",
			Message: "Product not found",
		})
		return
	}

	var req request.UpdateProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-024",
			Message: "Invalid request body",
		})
		return
	}

	if err := config.DB.Model(&product).Update("is_active", req.IsActive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "PRD-025",
			Message: "Failed to update product status",
		})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(product))
}

// toProductResponse converts a Product model to ProductResponse DTO
func toProductResponse(p models.Product) request.ProductResponse {
	resp := request.ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Slug:        p.Slug,
		SKU:         p.SKU,
		Description: p.Description,
		Price:       p.Price,
		SalePrice:   p.SalePrice,
		Thumbnail:   p.Thumbnail,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Category != nil {
		resp.Category = &request.ProductCategoryInfo{
			ID:   p.Category.ID,
			Name: p.Category.Name,
			Code: p.Category.Code,
		}
	}
	return resp
}
