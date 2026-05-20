package controllers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"demodiqit_api/config"
	"demodiqit_api/helpers/normalize"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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

	db := config.DB.Model(&models.Product{})

	if query.Keyword != "" {
		keyword := "%" + strings.ToLower(query.Keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ?", keyword, keyword)
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

	c.JSON(http.StatusOK, request.ProductListResponse{
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
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "PRD-005",
			Message: "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(product))
}

// CreateProduct handles POST /products
// Creates a new product. Slug/SKU are auto-cleaned by BeforeSave hook.
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-006",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	price, err := decimal.NewFromString(req.Price)
	if err != nil || price.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-007",
			Message: "Invalid price value",
		})
		return
	}

	var salePrice decimal.Decimal
	if req.SalePrice != "" {
		salePrice, err = decimal.NewFromString(req.SalePrice)
		if err != nil {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "PRD-008",
				Message: "Invalid sale price value",
			})
			return
		}
	}

	var product models.Product
	// Normalize slug and SKU from user input; BeforeSave hook will also clean before DB insert
	cleanSlug := normalize.Slug(req.Slug)
	if cleanSlug != "" && product.ExistsBySlug(config.DB, cleanSlug) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "PRD-026",
			Message: "Slug already exists. Please use a different product name or enter a custom slug.",
		})
		return
	}

	var cleanSKU string
	if req.SKU != "" {
		cleanSKU = normalize.SKU(req.SKU)
		if product.ExistsBySKU(config.DB, cleanSKU) {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "PRD-020",
				Message: "SKU already exists",
			})
			return
		}
	}

	product = models.Product{
		Name:        req.Name,
		Slug:        cleanSlug,
		SKU:         cleanSKU,
		Description: req.Description,
		Price:       price,
		SalePrice:   salePrice,
		Thumbnail:   req.Thumbnail,
		IsActive:    true,
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

	// Handle slug: normalize and validate uniqueness
	if req.Slug != "" {
		cleanSlug := normalize.Slug(req.Slug)
		if product.ExistsBySlug(config.DB, cleanSlug, product.ID) {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "PRD-018",
				Message: "Slug already exists",
			})
			return
		}
		product.Slug = cleanSlug
	}

	// Handle SKU: normalize and validate uniqueness
	if req.SKU != "" {
		cleanSKU := normalize.SKU(req.SKU)
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
	if req.Price != "" {
		price, err := decimal.NewFromString(req.Price)
		if err != nil || price.LessThanOrEqual(decimal.Zero) {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "PRD-013",
				Message: "Invalid price value",
			})
			return
		}
		product.Price = price
	}
	if req.SalePrice != "" {
		salePrice, err := decimal.NewFromString(req.SalePrice)
		if err != nil {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "PRD-014",
				Message: "Invalid sale price value",
			})
			return
		}
		product.SalePrice = salePrice
	}
	if req.Thumbnail != "" {
		product.Thumbnail = req.Thumbnail
	}

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
	return request.ProductResponse{
		ID:          p.ID,
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
}
