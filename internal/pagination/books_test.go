package pagination

import (
	"bookshop/internal/assert"
	"bookshop/internal/models"
	"fmt"
	"maps"
	"net/http"
	"testing"

	"github.com/govalues/decimal"
)

func Test_ParseBookPaginationQuery(t *testing.T) {
	tests := []struct {
		name                string
		query               string
		expectedMetaData    *MetaData
		expectedBookFilters *BookFilters
		expectedErr         string
	}{
		{
			name:  "Valid query parameters",
			query: "page=2&size=20&sort=created_at.desc&format=audiobook,e-book&category=1,2&min_price=10&max_price=100",
			expectedMetaData: &MetaData{
				PageNumber: 2,
				PageSize:   20,
				SortBy: map[string]string{
					"created_at": "desc",
				},
			},
			expectedBookFilters: &BookFilters{
				Formats:    []models.BookFormat{models.Audiobook, models.EBook},
				Categories: []int64{1, 2},
				MinPrice:   decimal.Ten,
				MaxPrice:   decimal.Hundred,
			},
			expectedErr: "",
		},
		{
			name:  "Missing query parameters (default values)",
			query: "",
			expectedMetaData: &MetaData{
				PageNumber: 0,
				PageSize:   10,
				SortBy:     nil,
			},
			expectedBookFilters: &BookFilters{
				Formats:    []models.BookFormat{models.Paperback},
				Categories: nil,
				MinPrice:   decimal.Zero,
				MaxPrice:   decimal.Zero,
			},
			expectedErr: "",
		},
		{
			name:                "Invalid format (invalid format type)",
			query:               "page=1&size=10&sort=created_at.desc&format=invalid_format",
			expectedMetaData:    nil,
			expectedBookFilters: nil,
			expectedErr:         "invalid book format parameter: invalid_format",
		},
		{
			name:                "Invalid category (invalid ID)",
			query:               "page=1&size=10&sort=created_at.desc&category=invalid_category",
			expectedMetaData:    nil,
			expectedBookFilters: nil,
			expectedErr:         "invalid book category ID invalid_category",
		},
		{
			name:                "Invalid min price (not a number)",
			query:               "page=1&size=10&sort=created_at.desc&min_price=abc",
			expectedMetaData:    nil,
			expectedBookFilters: nil,
			expectedErr:         "invalid min price abc",
		},
		{
			name:                "Invalid max price (not a number)",
			query:               "page=1&size=10&sort=created_at.desc&max_price=xyz",
			expectedMetaData:    nil,
			expectedBookFilters: nil,
			expectedErr:         "invalid max price xyz",
		},
		{
			name:                "Max price less than min price",
			query:               "page=1&size=10&sort=created_at.desc&min_price=50.0&max_price=30.0",
			expectedMetaData:    nil,
			expectedBookFilters: nil,
			expectedErr:         "max price must be greater than min price. min: 50.00, max: 30.00",
		},
		{
			name:  "Valid query with no filters",
			query: "page=1&size=10&sort=created_at.desc",
			expectedMetaData: &MetaData{
				PageNumber: 1,
				PageSize:   10,
				SortBy: map[string]string{
					"created_at": "desc",
				},
			},
			expectedBookFilters: &BookFilters{
				Formats:    []models.BookFormat{models.Paperback},
				Categories: nil,
				MinPrice:   decimal.Zero,
				MaxPrice:   decimal.Zero,
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/?"+tt.query, nil)
			assert.NoError(t, err)

			m, bf, err := ParseBookPaginationQuery(req)

			if tt.expectedErr != "" {
				assert.StringContains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, m.PageNumber, tt.expectedMetaData.PageNumber)
				assert.Equal(t, m.PageSize, tt.expectedMetaData.PageSize)
				assert.True(t, maps.Equal(m.SortBy, tt.expectedMetaData.SortBy))
				assert.Equal(t, len(bf.Formats), len(tt.expectedBookFilters.Formats))
				assert.Equal(t, len(bf.Categories), len(tt.expectedBookFilters.Categories))
				assert.True(t, bf.MinPrice.Equal(tt.expectedBookFilters.MinPrice))
				assert.True(t, bf.MaxPrice.Equal(tt.expectedBookFilters.MaxPrice))
			}
		})
	}
}

func Test_BuildFilterQuery(t *testing.T) {
	tests := []struct {
		name        string
		bookFilters BookFilters
		expected    string
	}{
		{
			name: "Filters with available field only",
			bookFilters: BookFilters{
				BookAvailability: models.Available,
			},
			expected: fmt.Sprintf("WHERE available = '%s'", models.Available),
		},
		{
			name: "Filters with formats",
			bookFilters: BookFilters{
				BookAvailability: models.Available,
				Formats:          []models.BookFormat{models.Audiobook, models.EBook},
			},
			expected: fmt.Sprintf("WHERE available = '%s' AND format IN ('%s', '%s')", models.Available, models.Audiobook, models.EBook),
		},
		{
			name: "Filters with categories",
			bookFilters: BookFilters{
				BookAvailability: models.Available,
				Categories:       []int64{1, 2},
			},
			expected: fmt.Sprintf("WHERE available = '%s' AND c.id IN (%d, %d)", models.Available, 1, 2),
		},
		{
			name: "Filters with price range",
			bookFilters: BookFilters{
				BookAvailability: models.Available,
				MinPrice:         decimal.Ten,
				MaxPrice:         decimal.Hundred,
			},
			expected: fmt.Sprintf("WHERE available = '%s' AND price >= %v AND price <= %v", models.Available, decimal.Ten, decimal.Hundred),
		},
		{
			name: "Filters with all fields",
			bookFilters: BookFilters{
				BookAvailability: models.Available,
				Formats:          []models.BookFormat{models.EBook, models.Paperback},
				Categories:       []int64{3, 4},
				MinPrice:         decimal.Ten,
				MaxPrice:         decimal.Hundred,
			},
			expected: fmt.Sprintf("WHERE available = '%s' AND format IN ('%s', '%s') AND c.id IN (%d, %d) AND price >= %v AND price <= %v", models.Available, models.EBook, models.Paperback, 3, 4, decimal.Ten, decimal.Hundred),
		},
		{
			name: "Filters with no formats or categories",
			bookFilters: BookFilters{
				BookAvailability: models.NotAvailable,
			},
			expected: fmt.Sprintf("WHERE available = '%s'", models.NotAvailable),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.bookFilters.BuildFilterQuery()
			assert.NormalizedStringsEqual(t, res, tt.expected)
		})
	}
}
