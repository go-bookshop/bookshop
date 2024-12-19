package doc

import (
	"bookshop/internal/models"
	"bookshop/internal/pagination"
)

const BooksTag = "Books"

func getBooksOperation() *Operation {
	op := NewOperation("Get Books", "Returns paginated books result", []string{BooksTag})
	op.AddQueryParameter(pagination.PageNumberParam, "integer", "Requested page number", false)
	op.AddQueryParameter(pagination.PageSizeParam, "integer", "Number of items per page", false)
	op.AddQueryParameter(pagination.SortParam, "string", "Sorting filters", false)
	op.AddQueryParameter(pagination.BookFormatParam, "string", "Book formats filter", false)
	op.AddQueryParameter(pagination.BookCategoryParam, "string", "Book categories filter", false)
	op.AddQueryParameter(pagination.BookMinPriceParam, "number", "Book's minimal price", false)
	op.AddQueryParameter(pagination.BookMinPriceParam, "number", "Book's maximal price", false)

	successSchema := NewSchema()

	successSchema.AddProperty("page", "integer", "Current page number", false)
	successSchema.AddProperty("size", "integer", "Number of items per page", false)
	successSchema.AddProperty("max_pages", "integer", "Maximum number of pages with a specified page size", false)

	refSchema := NewRefSchema("BookItem")
	successSchema.AddSchemaArrayProperty("data", "Paginated array of books", refSchema, false)

	requestErrorSchema := NewSchema()
	requestErrorSchema.AddProperty("errMsg", "string", "Specific error", false)

	op.AddResponseSchema(successSchema, "200", "application/json", "OK")
	op.AddResponseSchema(requestErrorSchema, "400", "application/json", "Bad request")

	return op
}

func getBookItemComponentSchema() *Schema {
	component := NewSchema()

	component.AddProperty("id", "integer", "Generated id of the book", false)
	component.AddProperty("title", "string", "Title of the book", false)
	component.AddProperty("synopsis", "string", "Synopsis of the book", false)
	component.AddProperty("avg_review", "number", "Average reviews of the book", false)
	component.AddSimpleArrayProperty("image_urls", "string", "Book images url", false)
	component.AddProperty("authors", "string", "Comma-separated list of authors of the book", false)
	component.AddProperty("categories", "string", "Comma-separated list of book's categories", false)

	propertyRefSchema := NewRefSchema("BookProperty")
	component.AddRefSchemaProperty("property", propertyRefSchema, false)

	return component
}

func getBookPropertySchema() *Schema {
	component := NewSchema()

	component.AddProperty("id", "integer", "Generated id of the book property", false)
	component.AddProperty("isbn", "string", "Unique book's ISBN", false)
	component.AddProperty("book_id", "integer", "Generated id of the related book", false)
	component.AddEnumProperty(
		"format",
		"string",
		"Format of the related book",
		false,
		string(models.Paperback),
		string(models.Hardcover),
		string(models.Audiobook),
		string(models.EBook),
	)
	component.AddProperty("language", "string", "Language of the related book", false)
	component.AddProperty("price", "number", "Price of the related book in current format", false)
	component.AddProperty("publisher", "string", "Publisher of the related book", false)
	component.AddProperty("target_audience", "string", "Target audience of the related book", false)
	component.AddProperty("illustrator", "string", "Illustrator of the related book", false)
	component.AddProperty("illustrations", "string", "Illustrations of the related book", false)
	component.AddProperty("page_number", "integer", "Number of pages of the related book", false)
	component.AddProperty("available", "string", "Availability of the related book in current format", false)
	component.AddEnumProperty(
		"available",
		"string",
		"Availability of the related book in current format",
		false,
		string(models.Available),
		string(models.NotAvailable),
		string(models.Upcoming),
	)
	component.AddProperty("published_at", "string", "Publish date of the related book", false)

	return component
}
