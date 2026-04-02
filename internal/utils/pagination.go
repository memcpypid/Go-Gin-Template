package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Sort   string `json:"sort"`
	SortBy string `json:"sort_by"`
}

func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	limit := 10
	page := 1
	sort := "desc"
	sortBy := "created_at"

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sortStr := c.Query("sort"); sortStr == "asc" || sortStr == "desc" {
		sort = sortStr
	}

	if sortByStr := c.Query("sort_by"); sortByStr != "" {
		sortBy = sortByStr
	}

	return Pagination{
		Limit:  limit,
		Page:   page,
		Sort:   sort,
		SortBy: sortBy,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
