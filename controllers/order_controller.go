package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"demodiqit_api/config"
	"demodiqit_api/helpers/context"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
)

// OrderController handles order-related HTTP requests
type OrderController struct {
	cfg *config.Config
}

// NewOrderController creates a new OrderController instance
func NewOrderController(cfg *config.Config) *OrderController {
	return &OrderController{cfg: cfg}
}

// ListOrders handles GET /orders
// Returns a paginated list of orders with optional filtering
func (oc *OrderController) ListOrders(c *gin.Context) {
	var query request.OrderListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-001",
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

	db := config.DB.Model(&models.Order{})

	// Filter by keyword (order number or customer name/email/phone)
	if query.Keyword != "" {
		keyword := "%" + strings.ToLower(query.Keyword) + "%"
		db = db.Where("LOWER(order_number) LIKE ? OR LOWER(customer_name) LIKE ? OR LOWER(customer_email) LIKE ? OR customer_phone LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	// Filter by status
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// Filter by date range
	if query.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", query.DateFrom); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if query.DateTo != "" {
		if t, err := time.Parse("2006-01-02", query.DateTo); err == nil {
			// Add 1 day to include the end date
			db = db.Where("created_at < ?", t.AddDate(0, 0, 1))
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-002",
			Message: "Failed to count orders",
		})
		return
	}

	var orders []models.Order
	if err := db.Preload("OrderItems").Order("created_at DESC").Offset(offset).Limit(query.Limit).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-003",
			Message: "Failed to fetch orders",
		})
		return
	}

	items := make([]request.OrderResponse, len(orders))
	for i, o := range orders {
		items[i] = toOrderResponse(o, false)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))

	c.JSON(http.StatusOK, request.OrderListResponse{
		Items:      items,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
	})
}

// GetOrder handles GET /orders/:id
// Returns details of a single order
func (oc *OrderController) GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-004",
			Message: "Invalid order ID",
		})
		return
	}

	var order models.Order
	if err := config.DB.Preload("OrderItems").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "ORD-005",
			Message: "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order, true))
}

// UpdateOrderStatus handles PATCH /orders/:id/status
// Updates only the status of an order
func (oc *OrderController) UpdateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-006",
			Message: "Invalid order ID",
		})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "ORD-007",
			Message: "Order not found",
		})
		return
	}

	var req request.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-008",
			Message: "Invalid request body",
		})
		return
	}

	// Validate status
	if !request.IsValidStatus(req.Status) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-009",
			Message: "Invalid status value. Must be: pending, processing, completed, or cancelled",
		})
		return
	}

	if err := config.DB.Model(&order).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-010",
			Message: "Failed to update order status",
		})
		return
	}

	// Reload with order items
	config.DB.Preload("OrderItems").First(&order, id)

	c.JSON(http.StatusOK, toOrderResponse(order, true))
}

// toOrderResponse converts an Order model to OrderResponse DTO
func toOrderResponse(o models.Order, includeItems bool) request.OrderResponse {
	resp := request.OrderResponse{
		ID:              o.ID,
		OrderNumber:     o.OrderNumber,
		Status:          o.Status,
		Subtotal:        o.Subtotal,
		ShippingFee:     o.ShippingFee,
		TotalAmount:     o.TotalAmount,
		CustomerName:    o.CustomerName,
		CustomerEmail:   o.CustomerEmail,
		CustomerPhone:   o.CustomerPhone,
		ShippingAddress: o.ShippingAddress,
		Notes:           o.Notes,
		ItemsCount:      len(o.OrderItems),
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}

	if includeItems {
		items := make([]request.OrderItemResponse, len(o.OrderItems))
		for i, item := range o.OrderItems {
			items[i] = request.OrderItemResponse{
				ID:               item.ID,
				ProductID:        item.ProductID,
				ProductName:      item.ProductName,
				ProductThumbnail: item.ProductThumbnail,
				Quantity:         item.Quantity,
				Price:            item.Price,
				Subtotal:         item.Subtotal,
			}
		}
		resp.OrderItems = items
	}

	return resp
}

// CreateOrder handles POST /orders
// Creates a new order with items
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-011",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate items exist
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-012",
			Message: "At least one item is required",
		})
		return
	}

	// Collect product IDs
	productIDs := make([]uint, len(req.Items))
	for i, item := range req.Items {
		productIDs[i] = item.ProductID
	}

	// Fetch all products in a single query
	var products []models.Product
	if err := config.DB.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-015",
			Message: "Failed to fetch products",
		})
		return
	}

	// Build product map for quick lookup
	productMap := make(map[uint]models.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	// Check if all products exist
	if len(productMap) != len(req.Items) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-013",
			Message: "One or more products not found",
		})
		return
	}

	// Generate order number
	orderNumber := generateOrderNumber()

	// Calculate subtotal from items
	var subtotal int64 = 0
	orderItems := make([]models.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		product := productMap[item.ProductID]
		itemSubtotal := int64(product.Price) * int64(item.Quantity)
		subtotal += itemSubtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:        &item.ProductID,
			ProductName:      product.Name,
			ProductThumbnail: product.Thumbnail,
			Quantity:         item.Quantity,
			Price:            int64(product.Price),
			Subtotal:         itemSubtotal,
		})
	}

	// Calculate total
	totalAmount := float64(subtotal) + req.ShippingFee

	userCtx := contextHelper.GetUserFromContext(c)
	userID := userCtx.ID

	order := models.Order{
		UserID:          &userID,
		OrderNumber:     orderNumber,
		Status:          request.OrderStatusPending,
		Subtotal:        float64(subtotal),
		ShippingFee:     req.ShippingFee,
		TotalAmount:     totalAmount,
		CustomerName:    req.CustomerName,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
		ShippingAddress: req.ShippingAddress,
		Notes:           req.Notes,
		OrderItems:      orderItems,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-014",
			Message: "Failed to create order",
		})
		return
	}

	config.DB.Preload("OrderItems").First(&order, order.ID)

	c.JSON(http.StatusCreated, respond.SuccessRespond{
		Message: "Order created successfully",
		Data:    toOrderResponse(order, true),
	})
}

// GuestCreateOrder handles POST /guest/orders
// Creates a new order without requiring authentication (for guest checkout)
func GuestCreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-011",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-012",
			Message: "At least one item is required",
		})
		return
	}

	productIDs := make([]uint, len(req.Items))
	for i, item := range req.Items {
		productIDs[i] = item.ProductID
	}

	var products []models.Product
	if err := config.DB.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-015",
			Message: "Failed to fetch products",
		})
		return
	}

	productMap := make(map[uint]models.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	if len(productMap) != len(req.Items) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-013",
			Message: "One or more products not found",
		})
		return
	}

	orderNumber := fmt.Sprintf("#ORD-%d", time.Now().UnixNano()/1000000%1000000)

	var subtotal int64 = 0
	orderItems := make([]models.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		product := productMap[item.ProductID]
		itemSubtotal := int64(product.Price) * int64(item.Quantity)
		subtotal += itemSubtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:        &item.ProductID,
			ProductName:      product.Name,
			ProductThumbnail: product.Thumbnail,
			Quantity:         item.Quantity,
			Price:            int64(product.Price),
			Subtotal:         itemSubtotal,
		})
	}

	totalAmount := float64(subtotal) + req.ShippingFee

	order := models.Order{
		OrderNumber:     orderNumber,
		Status:          request.OrderStatusPending,
		Subtotal:        float64(subtotal),
		ShippingFee:     req.ShippingFee,
		TotalAmount:     totalAmount,
		CustomerName:    req.CustomerName,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
		ShippingAddress: req.ShippingAddress,
		Notes:           req.Notes,
		OrderItems:      orderItems,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-014",
			Message: "Failed to create order",
		})
		return
	}

	config.DB.Preload("OrderItems").First(&order, order.ID)

	c.JSON(http.StatusCreated, toOrderResponse(order, true))
}

// generateOrderNumber creates a unique order number
func generateOrderNumber() string {
	return fmt.Sprintf("#ORD-%d", time.Now().UnixNano()/1000000%1000000)
}

// MyOrders handles GET /my-orders — returns orders for the authenticated user
func MyOrders(c *gin.Context) {
	userCtx := contextHelper.GetUserFromContext(c)
	fmt.Printf("[MyOrders] userCtx.ID=%d\n", userCtx.ID)
	if userCtx.ID == 0 {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "ORD-020",
			Message: "User not authenticated",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	if err := config.DB.Model(&models.Order{}).Where("user_id = ?", userCtx.ID).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-021",
			Message: "Failed to count orders",
		})
		return
	}

	var orders []models.Order
	if err := config.DB.Preload("OrderItems").
		Where("user_id = ?", userCtx.ID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-022",
			Message: "Failed to fetch orders",
		})
		return
	}

	items := make([]request.OrderResponse, len(orders))
	for i, o := range orders {
		items[i] = toOrderResponse(o, false)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Orders retrieved successfully",
		Data: respond.PaginatedData{
			Items:      items,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	})
}

// MyOrderDetail handles GET /my-orders/:id — returns a single order for the authenticated user
func MyOrderDetail(c *gin.Context) {
	userCtx := contextHelper.GetUserFromContext(c)
	if userCtx.ID == 0 {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "ORD-023",
			Message: "User not authenticated",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-024",
			Message: "Invalid order ID",
		})
		return
	}

	var order models.Order
	if err := config.DB.Preload("OrderItems").Where("id = ? AND user_id = ?", id, userCtx.ID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "ORD-025",
			Message: "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Order retrieved successfully",
		Data:    toOrderResponse(order, true),
	})
}

// CancelMyOrder handles POST /my-orders/:id/cancel — cancels the order if it is still pending
func CancelMyOrder(c *gin.Context) {
	userCtx := contextHelper.GetUserFromContext(c)
	if userCtx.ID == 0 {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "ORD-026",
			Message: "User not authenticated",
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "ORD-027",
			Message: "Invalid order ID",
		})
		return
	}

	var order models.Order
	if err := config.DB.Where("id = ? AND user_id = ?", id, userCtx.ID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorRespond{
			Code:    "ORD-028",
			Message: "Order not found",
		})
		return
	}

	if order.Status != request.OrderStatusPending {
		c.JSON(http.StatusForbidden, respond.ErrorRespond{
			Code:    "ORD-029",
			Message: "Only pending orders can be cancelled",
		})
		return
	}

	if err := config.DB.Model(&order).Update("status", request.OrderStatusCancelled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "ORD-030",
			Message: "Failed to cancel order",
		})
		return
	}

	config.DB.Preload("OrderItems").First(&order, id)

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Order cancelled successfully",
		Data:    toOrderResponse(order, true),
	})
}
