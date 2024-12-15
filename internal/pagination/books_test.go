package pagination

import (
	"bookshop/internal/models"
	"maps"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseBookPaginationQuery(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		expectedData *BookPaginationData
		expectedErr  string
	}{
		{
			name:  "Valid query parameters",
			query: "page=2&size=20&sort=created_at.desc&format=audiobook,e-book&category=1,2&min_price=10&max_price=50",
			expectedData: &BookPaginationData{
				PaginationData: &PaginationData{
					PageNumber: 2,
					PageSize:   20,
					SortBy: map[string]string{
						"created_at": "desc",
					},
				},
				Format:   []string{string(models.Audiobook), string(models.EBook)},
				Category: []int64{1, 2},
				MinPrice: 10,
				MaxPrice: 50,
			},
			expectedErr: "",
		},
		{
			name:  "Missing query parameters (default values)",
			query: "",
			expectedData: &BookPaginationData{
				PaginationData: &PaginationData{
					PageNumber: 0,
					PageSize:   10,
					SortBy:     nil,
				},
				Format:   []string{string(models.Paperback)},
				Category: nil,
				MinPrice: 0,
				MaxPrice: 0,
			},
			expectedErr: "",
		},
		{
			name:         "Invalid format (invalid format type)",
			query:        "page=1&size=10&sort=created_at.desc&format=invalid_format",
			expectedData: nil,
			expectedErr:  "invalid book format parameter: invalid_format",
		},
		{
			name:         "Invalid category (invalid ID)",
			query:        "page=1&size=10&sort=created_at.desc&category=invalid_category",
			expectedData: nil,
			expectedErr:  "invalid book category ID invalid_category",
		},
		{
			name:         "Invalid min price (not a number)",
			query:        "page=1&size=10&sort=created_at.desc&min_price=abc",
			expectedData: nil,
			expectedErr:  "invalid min price abc",
		},
		{
			name:         "Invalid max price (not a number)",
			query:        "page=1&size=10&sort=created_at.desc&max_price=xyz",
			expectedData: nil,
			expectedErr:  "invalid max price xyz",
		},
		{
			name:         "Max price less than min price",
			query:        "page=1&size=10&sort=created_at.desc&min_price=50&max_price=30",
			expectedData: nil,
			expectedErr:  "max price must be greater than min price. min: 50, max: 30",
		},
		{
			name:  "Valid query with no filters",
			query: "page=1&size=10&sort=created_at.desc",
			expectedData: &BookPaginationData{
				PaginationData: &PaginationData{
					PageNumber: 1,
					PageSize:   10,
					SortBy: map[string]string{
						"created_at": "desc",
					},
				},
				Format:   []string{string(models.Paperback)},
				Category: nil,
				MinPrice: 0,
				MaxPrice: 0,
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/?"+tt.query, nil)
			assert.NoError(t, err)

			p, err := ParseBookPaginationQuery(req)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, p.PageNumber, tt.expectedData.PageNumber)
				assert.Equal(t, p.PageSize, tt.expectedData.PageSize)
				assert.True(t, maps.Equal(p.SortBy, tt.expectedData.SortBy))
				assert.ElementsMatch(t, p.Format, tt.expectedData.Format)
				assert.ElementsMatch(t, p.Category, tt.expectedData.Category)
				assert.Equal(t, p.MinPrice, tt.expectedData.MinPrice)
				assert.Equal(t, p.MaxPrice, tt.expectedData.MaxPrice)
			}
		})
	}
}
