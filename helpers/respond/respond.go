package respond

type ErrorRespond struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessRespond struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
