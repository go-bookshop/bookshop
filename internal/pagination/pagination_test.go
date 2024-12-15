package pagination

import (
	"bookshop/internal/assert"
	"fmt"
	"maps"
	"net/http"
	"strings"
	"testing"
)

func Test_ParsePaginationData(t *testing.T) {
	pageSizeErr := fmt.Sprintf("failed to parse page size: page size must be between 1 and %d", maxPageSize)
	pageNumberErr := "failed to parse page number: page number must be positive"
	parseErr := "failed to parse"

	tests := []struct {
		name         string
		query        string
		expectedData *PaginationData
		expectedErr  string
	}{
		{
			name:  "Valid query parameters",
			query: "page=2&size=20&sort=created_at.desc",
			expectedData: &PaginationData{
				PageNumber: 2,
				PageSize:   20,
				SortBy: map[string]string{
					"created_at": "desc",
				},
			},
			expectedErr: "",
		},
		{
			name:  "Missing query parameters (default values)",
			query: "",
			expectedData: &PaginationData{
				PageNumber: 0,
				PageSize:   10,
			},
			expectedErr: "",
		},
		{
			name:         "Invalid page number (negative)",
			query:        "page=-1&size=10",
			expectedData: nil,
			expectedErr:  pageNumberErr,
		},
		{
			name:         "Invalid page size (too small)",
			query:        "page=1&size=0",
			expectedData: nil,
			expectedErr:  pageSizeErr,
		},
		{
			name:         "Invalid page size (too large)",
			query:        "page=1&size=200",
			expectedData: nil,
			expectedErr:  pageSizeErr,
		},
		{
			name:         "Invalid page number (not an integer)",
			query:        "page=abc&size=10",
			expectedData: nil,
			expectedErr:  parseErr,
		},
		{
			name:         "Invalid page size (not an integer)",
			query:        "page=1&size=abc",
			expectedData: nil,
			expectedErr:  parseErr,
		},
		{
			name:         "Invalid format (missing order)",
			query:        "sort=created_at",
			expectedData: nil,
			expectedErr:  "sort parameter must be in the format 'field.order'",
		},
		{
			name:         "Invalid sort order",
			query:        "sort=created_at.invalid",
			expectedData: nil,
			expectedErr:  "must be 'asc' or 'desc'",
		},
		{
			name:         "Invalid sort field",
			query:        "sort=invalid_field.asc",
			expectedData: nil,
			expectedErr:  "invalid sort field",
		},
		{
			name:         "Valid sort field with trailing invalid parameter",
			query:        "sort=created_at.desc,invalid_field.asc",
			expectedData: nil,
			expectedErr:  "invalid sort field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/?"+tt.query, nil)
			assert.NoError(t, err)

			p, err := ParsePaginationData(req, func(s string) bool { return !strings.Contains(s, "invalid") })

			if tt.expectedErr != "" {
				assert.StringContains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, p.PageNumber, tt.expectedData.PageNumber)
				assert.Equal(t, p.PageSize, tt.expectedData.PageSize)
				assert.True(t, maps.Equal(p.SortBy, tt.expectedData.SortBy))
			}
		})
	}
}

func Test_CalculateMaxPages(t *testing.T) {
	tests := []struct {
		maxItems int
		pageSize int
		expected int
	}{
		{
			maxItems: 10,
			pageSize: 5,
			expected: 2,
		},
		{
			maxItems: 10,
			pageSize: 10,
			expected: 1,
		},
		{
			maxItems: 10,
			pageSize: 15,
			expected: 1,
		},
		{
			maxItems: 10,
			pageSize: 3,
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run("Test_CalculateMaxPages", func(t *testing.T) {
			res := CalculateMaxPages(tt.maxItems, tt.pageSize)
			assert.Equal(t, res, tt.expected)
		})
	}
}

func Test_BuildSortingQuery(t *testing.T) {
	tests := []struct {
		name           string
		paginationData PaginationData
		expected       string
	}{
		{
			name: "Single sort by field",
			paginationData: PaginationData{
				SortBy: map[string]string{
					"created_at": "desc",
				},
			},
			expected: "ORDER BY created_at desc",
		},
		{
			name: "Nil sort by map",
			paginationData: PaginationData{
				SortBy: nil,
			},
			expected: "",
		},
		{
			name: "Empty sort by map",
			paginationData: PaginationData{
				SortBy: map[string]string{},
			},
			expected: "",
		},
		{
			name: "Multiple sort by fields",
			paginationData: PaginationData{
				SortBy: map[string]string{
					"created_at": "desc",
					"updated_at": "asc",
				},
			},
			expected: "ORDER BY created_at desc,updated_at asc",
		},
		{
			name: "Sort fields with extra comma",
			paginationData: PaginationData{
				SortBy: map[string]string{
					"created_at": "desc",
					"updated_at": "asc",
					"price":      "asc",
				},
			},
			expected: "ORDER BY created_at desc,updated_at asc,price asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.paginationData.BuildSortingQuery()
			assert.Equal(t, res, tt.expected)
		})
	}
}
