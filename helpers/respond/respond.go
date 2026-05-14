package respond

type ErrorRespond struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessRespond struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
