package response

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Success(message string, data interface{}) SuccessResponse {
	if data == nil {
		data = map[string]interface{}{}
	}
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func Error(message string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
	}
}
