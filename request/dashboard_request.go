package request

// Query parameters for dashboard revenue chart
type DashboardChartQuery struct {
	Period string `form:"period"` // "7d" | "1m" | "3m" | "custom"
	DateFrom string `form:"date_from"` // "YYYY-MM-DD" (for custom)
	DateTo   string `form:"date_to"`   // "YYYY-MM-DD" (for custom)
}

// DailyRevenue represents revenue for a single day
type DailyRevenue struct {
	Date   string  `json:"date" gorm:"column:date"`   // "YYYY-MM-DD"
	Amount float64 `json:"amount" gorm:"column:amount"` // Revenue in VND
}

// DashboardChartResponse is the response DTO for the dashboard revenue chart
type DashboardChartResponse struct {
	Items  []DailyRevenue `json:"items"`
	Total  float64        `json:"total"`  // Total revenue in the period
	Orders int64          `json:"orders"` // Total orders in the period
}
