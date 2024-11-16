package httputil

import (
	"bookshop/internal/assert"
	"fmt"
	"net/http"
	"testing"
)

func Test_ParsePaginationQuery(t *testing.T) {
	pageSizeErr := fmt.Sprintf("invalid page size. value should not be less than 1 and greater than %d", maxPageSize)
	pageNumberErr := "invalid page number. value should not be less than 0"
	parseErr := "failed to parse"

	tests := []struct {
		name         string
		query        string
		expectedData *PaginationData
		expectedErr  string
	}{
		{
			name:  "Valid query parameters",
			query: "page=2&size=20",
			expectedData: &PaginationData{
				PageNumber: 2,
				PageSize:   20,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/?"+tt.query, nil)
			assert.NoError(t, err)

			pagination, err := ParsePaginationQuery(req)

			if tt.expectedErr != "" {
				assert.NotNil(t, err)
				assert.StringContains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, pagination.PageNumber, tt.expectedData.PageNumber)
				assert.Equal(t, pagination.PageSize, tt.expectedData.PageSize)
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
