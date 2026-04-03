package response

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	validate *validator.Validate
)

func init() {
	// Initialize locales and translator
	enLocale := en.New()
	uni = ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")

	// Get the standard validator engine from Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		// Register default translations
		en_translations.RegisterDefaultTranslations(validate, trans)

		// Register custom name mapper to use JSON tags in error messages
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

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

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success    bool          `json:"success"`
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Errors     []ErrorDetail `json:"errors,omitempty"`
	Stack      string        `json:"stack,omitempty"`
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

func Error(statusCode int, message string, errs ...ErrorDetail) ErrorResponse {
	return ErrorResponse{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Errors:     errs,
	}
}

func ValidationError(err error) ErrorResponse {
	var errDetails []ErrorDetail

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			// Translate message automatically
			message := e.Translate(trans)

			errDetails = append(errDetails, ErrorDetail{
				Field:   fmt.Sprintf("body.%s", e.Field()),
				Message: message,
			})
		}
	} else {
		errDetails = append(errDetails, ErrorDetail{
			Field:   "body",
			Message: err.Error(),
		})
	}

	return ErrorResponse{
		Success:    false,
		StatusCode: 422,
		Message:    "Validation failed",
		Errors:     errDetails,
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
