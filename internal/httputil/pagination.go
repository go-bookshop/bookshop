package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const maxPageSize = 100

type PaginationData struct {
	PageNumber int
	PageSize   int
	SortBy     map[string]string
	FilterBy   map[string]string
}

func ParsePaginationQuery(r *http.Request) (*PaginationData, error) {
	qs := r.URL.Query()

	pagination := &PaginationData{
		PageNumber: 0,
		PageSize:   10,
	}

	pageNumber := qs.Get("page")
	if pageNumber != "" {
		pn, err := strconv.Atoi(pageNumber)

		if pn < 0 || err != nil {
			return nil, errors.New("invalid page number. value should not be less than 0")
		}

		pagination.PageNumber = pn
	}

	pageSize := qs.Get("size")
	if pageSize != "" {
		ps, err := strconv.Atoi(pageSize)

		if (ps < 1 || ps > maxPageSize) || err != nil {
			return nil, fmt.Errorf("invalid page size. value should not be less than 1 and greater than %d", maxPageSize)
		}

		pagination.PageSize = ps
	}

	return pagination, nil
}

func CalculateMaxPages(maxItems, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}

	return (maxItems + pageSize - 1) / pageSize
}
