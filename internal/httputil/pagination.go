package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const maxPageSize = 100

type PaginationData struct {
	PageNumber int
	PageSize   int
	SortBy     map[string]string
	FilterBy   map[string]string
}

func ParsePaginationQuery(r *http.Request, sortKeyValidator func(string) bool) (*PaginationData, error) {
	qs := r.URL.Query()

	pagination := &PaginationData{
		PageNumber: 0,
		PageSize:   10,
		SortBy:     make(map[string]string),
		FilterBy:   make(map[string]string),
	}

	pageNumber := qs.Get("page")
	if pageNumber != "" {
		pn, err := strconv.Atoi(pageNumber)

		if err != nil {
			return nil, fmt.Errorf("failed to parse page: %w", err)
		}
		if pn < 0 {
			return nil, errors.New("invalid page number. value should not be less than 0")
		}

		pagination.PageNumber = pn
	}

	pageSize := qs.Get("size")
	if pageSize != "" {
		ps, err := strconv.Atoi(pageSize)

		if err != nil {
			return nil, fmt.Errorf("failed to parse size: %w", err)
		}
		if ps < 1 || ps > maxPageSize {
			return nil, fmt.Errorf("invalid page size. value should not be less than 1 and greater than %d", maxPageSize)
		}

		pagination.PageSize = ps
	}

	sort := qs.Get("sort")
	if sort != "" {
		sortMap, err := parseSortQueryParameter(sort, sortKeyValidator)

		if err != nil {
			return nil, fmt.Errorf("invalid sort parameter. %w", err)
		}

		pagination.SortBy = sortMap
	}

	return pagination, nil
}

func parseSortQueryParameter(sort string, sortKeyValidator func(string) bool) (map[string]string, error) {
	sortMap := make(map[string]string)

	sortKeys := strings.Split(sort, "&")
	for _, sk := range sortKeys {
		parts := strings.Split(sk, ".")
		if len(parts) != 2 {
			return nil, errors.New("sort parameter must be in the format 'field.order'")
		}

		f, o := parts[0], parts[1]
		if o != "asc" && o != "desc" {
			return nil, fmt.Errorf("invalid sort order '%s'. must be 'asc' or 'desc'", o)
		}

		if ok := sortKeyValidator(f); !ok {
			return nil, fmt.Errorf("invalid sort field '%s'", f)
		}

		sortMap[f] = o
	}

	return sortMap, nil
}

func CalculateMaxPages(maxItems, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}

	return (maxItems + pageSize - 1) / pageSize
}

func (p *PaginationData) BuildSortingQuery() string {
	if p.SortBy == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("ORDER BY ")

	for k, v := range p.SortBy {
		sb.WriteString(fmt.Sprintf("%s %s,", k, v))
	}

	return strings.TrimSuffix(sb.String(), ",")
}
