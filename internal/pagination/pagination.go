package pagination

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	PageNumberParam = "page"
	PageSizeParam   = "size"
	SortParam       = "sort"

	defaultPageSize = 10
	maxPageSize     = 100
)

type MetaData struct {
	PageNumber int
	PageSize   int
	SortBy     map[string]string
}

func ParsePaginationMetaData(r *http.Request, sortKeyValidator func(string) bool) (*MetaData, error) {
	qs := r.URL.Query()

	pageNumber, err := parsePageNumber(qs.Get(PageNumberParam))
	if err != nil {
		return nil, fmt.Errorf("failed to parse page number: %w", err)
	}

	pageSize, err := parsePageSize(qs.Get(PageSizeParam))
	if err != nil {
		return nil, fmt.Errorf("failed to parse page size: %w", err)
	}

	sortBy, err := parseSortQuery(qs.Get(SortParam), sortKeyValidator)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sort parameter: %w", err)
	}

	return &MetaData{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		SortBy:     sortBy,
	}, nil
}

func parsePageNumber(pageNumberParam string) (int, error) {
	if pageNumberParam == "" {
		return 0, nil
	}

	pageNumber, err := strconv.Atoi(pageNumberParam)
	if err != nil {
		return 0, fmt.Errorf("page must be a valid integer, got '%s'", pageNumberParam)
	}

	if pageNumber < 0 {
		return 0, errors.New("page number must be positive")
	}

	return pageNumber, nil
}

func parsePageSize(pageSizeParam string) (int, error) {
	if pageSizeParam == "" {
		return defaultPageSize, nil
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil {
		return 0, fmt.Errorf("size must be a valid integer, got '%s'", pageSizeParam)
	}

	if pageSize < 1 || pageSize > maxPageSize {
		return 0, fmt.Errorf("page size must be between 1 and %d, got %d", maxPageSize, pageSize)
	}

	return pageSize, nil
}

func parseSortQuery(sortParam string, sortKeyValidator func(string) bool) (map[string]string, error) {
	if sortParam == "" {
		return nil, nil
	}

	sortMap := make(map[string]string)
	sortKeys := strings.Split(sortParam, ",")

	for _, sortKey := range sortKeys {
		field, order, err := parseSortField(sortKey)
		if err != nil {
			return nil, fmt.Errorf("invalid sort parameter '%s': %w", sortKey, err)
		}

		if !sortKeyValidator(field) {
			return nil, fmt.Errorf("invalid sort field '%s'", field)
		}

		sortMap[field] = order
	}

	return sortMap, nil
}

func parseSortField(sortField string) (string, string, error) {
	parts := strings.Split(sortField, ".")
	if len(parts) != 2 {
		return "", "", errors.New("sort parameter must be in the format 'field.order'")
	}

	field, order := parts[0], parts[1]
	if order != "asc" && order != "desc" {
		return "", "", fmt.Errorf("invalid sort order '%s'. must be 'asc' or 'desc'", order)
	}

	return field, order, nil
}

func CalculateMaxPages(maxItems, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}

	return (maxItems + pageSize - 1) / pageSize
}

func (p MetaData) BuildSortingQuery() string {
	if len(p.SortBy) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("ORDER BY ")

	for sortKey, sortOrder := range p.SortBy {
		sb.WriteString(fmt.Sprintf("%s %s,", sortKey, sortOrder))
	}

	return strings.TrimSuffix(sb.String(), ",")
}
