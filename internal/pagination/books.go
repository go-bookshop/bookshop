package pagination

import (
	"bookshop/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/govalues/decimal"
)

const (
	BookFormatParam   = "format"
	BookCategoryParam = "category"
	BookMinPriceParam = "min_price"
	BookMaxPriceParam = "max_price"
)

type BookFilters struct {
	BookAvailability models.BookAvailability
	Formats          []models.BookFormat
	Categories       []int64
	MinPrice         decimal.Decimal
	MaxPrice         decimal.Decimal
}

func ParseBookPaginationQuery(r *http.Request) (*MetaData, *BookFilters, error) {
	qs := r.URL.Query()

	paginationData, err := ParsePaginationMetaData(r, bookItemSortKeyValidator)
	if err != nil {
		return nil, nil, err
	}

	filters := &BookFilters{
		BookAvailability: models.Available,
		Formats:          []models.BookFormat{models.Paperback},
	}

	if formats, err := parseAndValidateFormats(qs.Get(BookFormatParam)); err != nil {
		return nil, nil, err
	} else if len(formats) > 0 {
		filters.Formats = formats
	}

	categories, err := parseCategories(qs.Get(BookCategoryParam))
	if err != nil {
		return nil, nil, err
	}
	filters.Categories = categories

	minPrice, maxPrice, err := parseMinMaxPrice(qs.Get(BookMinPriceParam), qs.Get(BookMaxPriceParam))
	if err != nil {
		return nil, nil, err
	}
	filters.MinPrice = minPrice
	filters.MaxPrice = maxPrice

	return paginationData, filters, nil
}

func parseAndValidateFormats(formatParam string) ([]models.BookFormat, error) {
	if formatParam == "" {
		return nil, nil
	}

	var formatFilters []models.BookFormat
	formats := strings.Split(formatParam, ",")

	for _, f := range formats {
		bookFormat := models.BookFormat(f)
		if isValid := bookFormat.Validate(); !isValid {
			return nil, fmt.Errorf("invalid book format parameter: %v", f)
		}
		formatFilters = append(formatFilters, bookFormat)
	}

	return formatFilters, nil
}

func parseCategories(categoryParam string) ([]int64, error) {
	if categoryParam == "" {
		return nil, nil
	}

	var categoryIDs []int64
	categories := strings.Split(categoryParam, ",")

	for _, c := range categories {
		categoryID, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid book category ID %v: %w", c, err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	return categoryIDs, nil
}

func parseMinMaxPrice(minPriceParam, maxPriceParam string) (minPrice decimal.Decimal, maxPrice decimal.Decimal, err error) {
	if minPriceParam != "" {
		minPrice, err = decimal.ParseExact(minPriceParam, 5)
		if err != nil {
			return decimal.Zero, decimal.Zero, fmt.Errorf("invalid min price %v: %w", minPriceParam, err)
		}
	}

	if maxPriceParam != "" {
		maxPrice, err = decimal.ParseExact(maxPriceParam, 5)
		if err != nil {
			return decimal.Zero, decimal.Zero, fmt.Errorf("invalid max price %v: %w", maxPriceParam, err)
		}
		if minPrice.Cmp(decimal.Zero) > 0 && maxPrice.Cmp(minPrice) < 1 {
			return decimal.Zero, decimal.Zero, fmt.Errorf("max price must be greater than min price. min: %v, max: %v", minPrice.Rescale(2), maxPrice.Rescale(2))
		}
	}

	return minPrice, maxPrice, nil
}

func bookItemSortKeyValidator(key string) bool {
	switch key {
	case "avg_review", "created_at", "price":
		return true
	}
	return false
}

func (bf BookFilters) BuildFilterQuery() string {
	var sb strings.Builder
	sb.WriteString("WHERE ")
	sb.WriteString(fmt.Sprintf("\n\tavailable = '%s' AND ", bf.BookAvailability))

	if len(bf.Formats) > 0 {
		sb.WriteString("\n\tformat IN (")
		for i, format := range bf.Formats {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("'%s'", format))
		}
		sb.WriteString(")")
		sb.WriteString("\n AND ")
	}

	if len(bf.Categories) > 0 {
		sb.WriteString("\n\tc.id IN (")
		for i, categoryID := range bf.Categories {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%v", categoryID))
		}
		sb.WriteString(")")
		sb.WriteString("\n AND ")
	}

	if bf.MinPrice.Cmp(decimal.Zero) > 0 {
		sb.WriteString(fmt.Sprintf("\n\tprice >= %v", bf.MinPrice))
		sb.WriteString("\n AND ")
	}

	if bf.MaxPrice.Cmp(decimal.Zero) > 0 {
		sb.WriteString(fmt.Sprintf("\n\tprice <= %v", bf.MaxPrice))
	}

	return strings.TrimSuffix(strings.TrimSpace(sb.String()), "AND")
}
