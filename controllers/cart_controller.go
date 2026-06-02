package controllers

import (
	"net/http"

	"demodiqit_api/config"
	"demodiqit_api/helpers/context"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/models"

	"github.com/gin-gonic/gin"
)

type cartController struct{}

func NewCartController() *cartController {
	return &cartController{}
}

// GetCart returns the user's cart
func (cc *cartController) GetCart(c *gin.Context) {
	userCtx := contextHelper.GetUserFromContext(c)
	if userCtx.ID == 0 {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "CART-001",
			Message: "User not authenticated",
		})
		return
	}

	var cart models.Cart
	err := config.DB.Preload("CartItems").
		Where("user_id = ?", userCtx.ID).
		First(&cart).Error

	if err != nil {
		c.JSON(http.StatusOK, respond.SuccessRespond{
			Message: "Cart retrieved successfully",
			Data:    map[string]interface{}{"cart_items": []interface{}{}, "total_items": 0},
		})
		return
	}

	totalItems := 0
	for _, item := range cart.CartItems {
		totalItems += item.Quantity
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Cart retrieved successfully",
		Data: map[string]interface{}{
			"id":           cart.ID,
			"cart_items":   cart.CartItems,
			"total_items":  totalItems,
		},
	})
}

// SaveCart saves or updates the entire cart (replaces current cart with new items)
func (cc *cartController) SaveCart(c *gin.Context) {
	userCtx := contextHelper.GetUserFromContext(c)
	if userCtx.ID == 0 {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "CART-002",
			Message: "User not authenticated",
		})
		return
	}

	var request struct {
		CartItems []struct {
			ProductID   uint    `json:"product_id" binding:"required"`
			ProductName string  `json:"product_name" binding:"required"`
			Thumbnail   string  `json:"thumbnail"`
			Price       float64 `json:"price" binding:"required"`
			Quantity    int     `json:"quantity" binding:"required,min=1"`
			Description string  `json:"description"`
		} `json:"cart_items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "CART-003",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	var cart models.Cart
	err := config.DB.Where("user_id = ?", userCtx.ID).First(&cart).Error
	if err != nil {
		cart = models.Cart{UserID: userCtx.ID}
		if err := config.DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
				Code:    "CART-004",
				Message: "Failed to create cart",
			})
			return
		}
	}

	// Delete existing cart items
	config.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})

	// Insert new cart items
	for _, item := range request.CartItems {
		cartItem := models.CartItem{
			CartID:      cart.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Thumbnail:   item.Thumbnail,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Description: item.Description,
		}
		if err := config.DB.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
				Code:    "CART-005",
				Message: "Failed to save cart items",
			})
			return
		}
	}

	// Reload cart with items
	config.DB.Preload("CartItems").First(&cart, cart.ID)

	totalItems := 0
	for _, item := range cart.CartItems {
		totalItems += item.Quantity
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Cart saved successfully",
		Data: map[string]interface{}{
			"id":          cart.ID,
			"cart_items":  cart.CartItems,
			"total_items": totalItems,
		},
	})
}
