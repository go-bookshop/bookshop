package pagination

import (
	"bookshop/internal/models"
	"maps"
	"net/http"
	"testing"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
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
				Format:   []models.BookFormat{models.Audiobook, models.EBook},
				Category: []int64{1, 2},
				MinPrice: decimal.Ten,
				MaxPrice: decimal.Hundred,
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
				Format:   []models.BookFormat{models.Paperback},
				Category: nil,
				MinPrice: decimal.Zero,
				MaxPrice: decimal.Zero,
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
				Format:   []models.BookFormat{models.Paperback},
				Category: nil,
				MinPrice: decimal.Zero,
				MaxPrice: decimal.Zero,
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
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, m.PageNumber, tt.expectedMetaData.PageNumber)
				assert.Equal(t, m.PageSize, tt.expectedMetaData.PageSize)
				assert.True(t, maps.Equal(m.SortBy, tt.expectedMetaData.SortBy))
				assert.ElementsMatch(t, bf.Format, tt.expectedBookFilters.Format)
				assert.ElementsMatch(t, bf.Category, tt.expectedBookFilters.Category)
				assert.True(t, bf.MinPrice.Equal(tt.expectedBookFilters.MinPrice))
				assert.True(t, bf.MaxPrice.Equal(tt.expectedBookFilters.MaxPrice))
			}
		})
	}
}
