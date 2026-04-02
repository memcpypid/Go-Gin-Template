package response

type PaginationMeta struct {
	Total       int64 `json:"total"`
	Limit       int   `json:"limit"`
	Page        int   `json:"page"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

type SuccessResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
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

func SuccessWithPagination(message string, data interface{}, total int64, limit, page int) SuccessResponse {
	if data == nil {
		data = []interface{}{}
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	hasNext := page < totalPages
	hasPrevious := page > 1

	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &PaginationMeta{
			Total:       total,
			Limit:       limit,
			Page:        page,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrevious: hasPrevious,
		},
	}
}
