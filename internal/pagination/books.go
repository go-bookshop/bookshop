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
	Availability string
	Format       []models.BookFormat
	Category     []int64
	MinPrice     decimal.Decimal
	MaxPrice     decimal.Decimal
}

func ParseBookPaginationQuery(r *http.Request) (*MetaData, *BookFilters, error) {
	qs := r.URL.Query()

	paginationData, err := ParsePaginationMetaData(r, bookItemSortKeyValidator)
	if err != nil {
		return nil, nil, err
	}

	filters := &BookFilters{
		Availability: string(models.Available),
		Format:       []models.BookFormat{models.Paperback},
	}

	if formats, err := parseAndValidateFormats(qs.Get(BookFormatParam)); err != nil {
		return nil, nil, err
	} else if len(formats) > 0 {
		filters.Format = formats
	}

	categories, err := parseCategories(qs.Get(BookCategoryParam))
	if err != nil {
		return nil, nil, err
	}
	filters.Category = categories

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

	var formatsMap []models.BookFormat
	formats := strings.Split(formatParam, ",")

	for _, f := range formats {
		bookFormat := models.BookFormat(f)
		switch bookFormat {
		case models.Audiobook, models.EBook, models.Hardcover, models.Paperback:
			formatsMap = append(formatsMap, bookFormat)
		default:
			return nil, fmt.Errorf("invalid book format parameter: %v", f)
		}
	}

	return formatsMap, nil
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

func parseMinMaxPrice(minPriceParam, maxPriceParam string) (decimal.Decimal, decimal.Decimal, error) {
	var minPrice, maxPrice decimal.Decimal
	var err error

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
	if key != "avg_review" &&
		key != "created_at" &&
		key != "price" {
		return false
	}

	return true
}

func (bf *BookFilters) BuildFilterQuery() string {
	var sb strings.Builder
	sb.WriteString("WHERE ")
	sb.WriteString(fmt.Sprintf("\n\tavailable = '%s' AND", bf.Availability))

	if len(bf.Format) > 0 {
		sb.WriteString("\n\tformat IN (")
		for i, format := range bf.Format {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("'%s'", format))
		}
		sb.WriteString(")")
		sb.WriteString("\nAND")
	}

	if len(bf.Category) > 0 {
		sb.WriteString("\n\tc.id IN (")
		for i, categoryID := range bf.Category {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("'%v'", categoryID))
		}
		sb.WriteString(")")
		sb.WriteString("\nAND")
	}

	if bf.MinPrice.Cmp(decimal.Zero) > 0 {
		sb.WriteString(fmt.Sprintf("\n\tprice >= %v", bf.MinPrice))
		sb.WriteString("\nAND")
	}

	if bf.MaxPrice.Cmp(decimal.Zero) > 0 {
		sb.WriteString(fmt.Sprintf("\n\tprice <= %v", bf.MaxPrice))
	}

	return strings.TrimSuffix(sb.String(), "AND")
}
