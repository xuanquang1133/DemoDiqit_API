package controllers

import (
	"fmt"
	"net/http"
	"time"

	"demodiqit_api/config"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardStats holds the aggregated stats shown on the CMS dashboard
type DashboardStats struct {
	TotalUsers       int64   `json:"total_users"`
	TotalProducts    int64   `json:"total_products"`
	CompletedOrders  int64   `json:"completed_orders"`
	TotalRevenue     float64 `json:"total_revenue"`
}

// MonthlyRevenue represents revenue for a single month
type MonthlyRevenue struct {
	Month  string  `json:"month"`  // e.g. "2025-01"
	Amount float64 `json:"amount"` // Revenue in VND (stored as float64)
}

// RecentOrderItem holds minimal info for the recent orders table
type RecentOrderItem struct {
	OrderNumber  string  `json:"order_number"`
	CustomerName string  `json:"customer_name"`
	ProductName  string  `json:"product_name"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status"`
}

// DashboardData holds all data for the dashboard page
type DashboardData struct {
	Stats          DashboardStats    `json:"stats"`
	RecentOrders   []RecentOrderItem `json:"recent_orders"`
	MonthlyRevenue []MonthlyRevenue  `json:"monthly_revenue"`
}

// GetDashboardStats handles GET /dashboard
// Returns aggregated stats, recent orders, and monthly revenue for the CMS dashboard
func GetDashboardStats(c *gin.Context) {
	// 1. Count total users
	var totalUsers int64
	if err := config.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "DASH-001",
			Message: "Failed to count users",
		})
		return
	}

	// 2. Count total active products
	var totalProducts int64
	if err := config.DB.Model(&models.Product{}).Where("is_active = ?", true).Count(&totalProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "DASH-002",
			Message: "Failed to count products",
		})
		return
	}

	// 3. Count completed orders only
	var completedOrders int64
	if err := config.DB.Model(&models.Order{}).
		Where("status = ?", "completed").
		Count(&completedOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "DASH-003",
			Message: "Failed to count completed orders",
		})
		return
	}

	// 4. Sum total revenue (completed orders only)
	var totalRevenue float64
	if err := config.DB.Model(&models.Order{}).
		Where("status = ?", "completed").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "DASH-004",
			Message: "Failed to calculate revenue",
		})
		return
	}

	// 5. Recent orders (last 5, with first product name)
	var orders []models.Order
	if err := config.DB.Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
		return db.Limit(1)
	}).
		Order("created_at DESC").
		Limit(5).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "DASH-005",
			Message: "Failed to fetch recent orders",
		})
		return
	}

	recentOrders := make([]RecentOrderItem, len(orders))
	for i, o := range orders {
		productName := ""
		if len(o.OrderItems) > 0 {
			productName = o.OrderItems[0].ProductName
		}
		recentOrders[i] = RecentOrderItem{
			OrderNumber:  o.OrderNumber,
			CustomerName: o.CustomerName,
			ProductName:  productName,
			TotalAmount:  o.TotalAmount,
			Status:       o.Status,
		}
	}

	// 6. Monthly revenue (last 12 months, completed orders only)
	var revenueResults []struct {
		Month  string
		Amount float64
	}
	if err := config.DB.Model(&models.Order{}).
		Select("TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM') AS month, COALESCE(SUM(total_amount), 0) AS amount").
		Where("status = ?", "completed").
		Group("DATE_TRUNC('month', created_at)").
		Order("month DESC").
		Limit(12).
		Scan(&revenueResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "DASH-006",
			Message: "Failed to calculate monthly revenue",
		})
		return
	}

	// Reverse so months go oldest → newest for chart rendering
	monthlyRevenue := make([]MonthlyRevenue, len(revenueResults))
	for i, r := range revenueResults {
		monthlyRevenue[len(revenueResults)-1-i] = MonthlyRevenue{
			Month:  r.Month,
			Amount: r.Amount,
		}
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "Success",
		Data: DashboardData{
			Stats: DashboardStats{
				TotalUsers:       totalUsers,
				TotalProducts:    totalProducts,
				CompletedOrders:   completedOrders,
				TotalRevenue:     totalRevenue,
			},
			RecentOrders:   recentOrders,
			MonthlyRevenue: monthlyRevenue,
		},
	})
}

// GetDashboardChart handles GET /dashboard/v2/chart
// Returns daily revenue data based on a time period filter
func GetDashboardChart(c *gin.Context) {
	var query request.DashboardChartQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "DASH-010",
			Message: "Invalid query parameters",
		})
		return
	}

	now := time.Now()
	var startDate, endDate time.Time

	switch query.Period {
	case "7d":
		startDate = now.AddDate(0, 0, -6)
		endDate = now
	case "1m":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case "3m":
		startDate = now.AddDate(0, -3, 0)
		endDate = now
	case "custom":
		if query.DateFrom == "" || query.DateTo == "" {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "DASH-011",
				Message: "Custom period requires date_from and date_to",
			})
			return
		}
		var err error
		startDate, err = time.Parse("2006-01-02", query.DateFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "DASH-012",
				Message: "Invalid date_from format. Expected YYYY-MM-DD",
			})
			return
		}
		endDate, err = time.Parse("2006-01-02", query.DateTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, respond.ErrorRespond{
				Code:    "DASH-013",
				Message: "Invalid date_to format. Expected YYYY-MM-DD",
			})
			return
		}
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	default:
		startDate = now.AddDate(0, 0, -29)
		endDate = now
	}

	// Validate date range: end must not be before start
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "DASH-015",
			Message: "date_to must be on or after date_from",
		})
		return
	}

	daysDiff := int(endDate.Sub(startDate).Hours() / 24)
	if daysDiff > 365 {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Code:    "DASH-016",
			Message: "Date range cannot exceed 365 days",
		})
		return
	}

	var revenueResults []request.DailyRevenue
	if err := config.DB.Raw(`
		SELECT
			TO_CHAR(DATE(created_at AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD') AS date,
			COALESCE(SUM(total_amount), 0) AS amount
		FROM orders
		WHERE status = 'completed'
		  AND created_at >= ?
		  AND created_at <= ?
		GROUP BY 1
		ORDER BY 1 ASC`,
		startDate, endDate,
	).Scan(&revenueResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Code:    "DASH-014",
			Message: "Failed to calculate revenue chart",
		})
		return
	}

	revenueMap := make(map[string]float64)
	for _, r := range revenueResults {
		revenueMap[r.Date] = r.Amount
	}

	items := make([]request.DailyRevenue, 0, daysDiff+1)
	for i := 0; i <= daysDiff; i++ {
		day := startDate.AddDate(0, 0, i)
		dateStr := day.Format("2006-01-02")
		items = append(items, request.DailyRevenue{
			Date:   dateStr,
			Amount: revenueMap[dateStr],
		})
	}

	var totalRevenue float64
	var totalOrders int64

	config.DB.Model(&models.Order{}).
		Where("status = ?", "completed").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Select("COALESCE(SUM(total_amount), 0), COUNT(*)").
		Row().Scan(&totalRevenue, &totalOrders)

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: fmt.Sprintf("Revenue from %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		Data: request.DashboardChartResponse{
			Items:  items,
			Total:  totalRevenue,
			Orders: totalOrders,
		},
	})
}
