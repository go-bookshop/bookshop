package pagination

import (
	"bookshop/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	formatParam   = "format"
	categoryParam = "category"
	minPriceParam = "min_price"
	maxPriceParam = "max_price"
)

type BookPaginationData struct {
	*PaginationData
	Format   []string
	Category []int64
	MinPrice float64
	MaxPrice float64
}

func ParseBookPaginationQuery(r *http.Request) (*BookPaginationData, error) {
	qs := r.URL.Query()

	p, err := ParsePaginationData(r, bookItemSortKeyValidator)
	if err != nil {
		return nil, err
	}

	bp := &BookPaginationData{PaginationData: p}

	if formats, err := parseAndValidateFormats(qs.Get(formatParam)); err != nil {
		return nil, err
	} else if len(formats) > 0 {
		bp.Format = formats
	} else {
		bp.Format = []string{string(models.Paperback)}
	}

	if categories, err := parseCategories(qs.Get(categoryParam)); err != nil {
		return nil, err
	} else {
		bp.Category = categories
	}

	if minPrice, maxPrice, err := parseMinMaxPrice(qs.Get(minPriceParam), qs.Get(maxPriceParam)); err != nil {
		return nil, err
	} else {
		bp.MinPrice = minPrice
		bp.MaxPrice = maxPrice
	}

	return bp, nil
}

func parseAndValidateFormats(formatParam string) ([]string, error) {
	if formatParam == "" {
		return nil, nil
	}

	formats := strings.Split(formatParam, ",")
	validFormats := []string{}

	for _, f := range formats {
		switch models.BookFormat(f) {
		case models.Audiobook, models.EBook, models.Hardcover, models.Paperback:
			validFormats = append(validFormats, f)
		default:
			return nil, fmt.Errorf("invalid book format parameter: %v", f)
		}
	}

	return validFormats, nil
}

func parseCategories(categoryParam string) ([]int64, error) {
	if categoryParam == "" {
		return nil, nil
	}

	categories := strings.Split(categoryParam, ",")
	categoryIDs := []int64{}

	for _, c := range categories {
		categoryID, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid book category ID %v: %w", c, err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	return categoryIDs, nil
}

func parseMinMaxPrice(minPriceParam, maxPriceParam string) (float64, float64, error) {
	var minPrice, maxPrice float64
	var err error

	if minPriceParam != "" {
		minPrice, err = strconv.ParseFloat(minPriceParam, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid min price %v: %w", minPriceParam, err)
		}
	}

	if maxPriceParam != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceParam, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid max price %v: %w", maxPriceParam, err)
		}
		if minPrice > 0 && maxPrice <= minPrice {
			return 0, 0, fmt.Errorf("max price must be greater than min price. min: %v, max: %v", minPrice, maxPrice)
		}
	}

	return minPrice, maxPrice, nil
}

func bookItemSortKeyValidator(key string) bool {
	if key != "avg_review" &&
		key != "created_at" &&
		key != "price" {
		return false
	}

	return true
}

func (bp *BookPaginationData) BuildFilterQuery() string {
	var sb strings.Builder
	sb.WriteString("WHERE ")
	sb.WriteString(fmt.Sprintf("available = '%s'", models.Available))

	if len(bp.Format) > 0 {
		sb.WriteString("\nformat IN (")
		for i, format := range bp.Format {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("'%s'", format))
		}
		sb.WriteString(")")
		sb.WriteString("\nAND")
	}

	if len(bp.Category) > 0 {
		sb.WriteString("\nc.id IN (")
		for i, categoryID := range bp.Category {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("'%v'", categoryID))
		}
		sb.WriteString(")")
		sb.WriteString("\nAND")
	}

	if bp.MinPrice > 0 {
		sb.WriteString(fmt.Sprintf("\nprice >= %v", bp.MinPrice))
		sb.WriteString("\nAND")
	}

	if bp.MaxPrice > 0 {
		sb.WriteString(fmt.Sprintf("\nprice <= %v", bp.MaxPrice))
	}

	return strings.TrimSuffix(sb.String(), "AND")
}
