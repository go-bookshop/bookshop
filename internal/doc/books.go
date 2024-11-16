package doc

const BooksTag = "Books"

func getBooksOperation() *Operation {
	op := NewOperation("Get Books", "Returns paginated books result", []string{BooksTag})
	op.AddQueryParameter("page", "integer", "Requested page number", false)
	op.AddQueryParameter("size", "integer", "Number of items per page", false)

	reqSchema := NoBodySchema

	successSchema := NewSchema()

	successSchema.AddProperty("page", "integer", "Current page number", true)
	successSchema.AddProperty("size", "integer", "Number of items per page", true)
	successSchema.AddProperty("max_pages", "integer", "Maximum number of pages with a specified page size", true)

	refSchema := NewRefSchema("BookItem")
	successSchema.AddSchemaArrayProperty("data", "Paginated array of books", refSchema, true)

	requestErrorSchema := NewSchema()
	requestErrorSchema.AddProperty("errMsg", "string", "Specific error", true)

	op.AddRequestSchema(reqSchema, "application/json", true)
	op.AddResponseSchema(successSchema, "200", "application/json", "OK")
	op.AddResponseSchema(requestErrorSchema, "400", "application/json", "Bad request")

	return op
}

func getBookItemComponentSchema() *Schema {
	component := NewSchema()

	component.AddProperty("id", "integer", "Generated id of the book", true)
	component.AddProperty("title", "string", "Title of the book", true)
	component.AddProperty("synopsis", "string", "Synopsis of the book", true)
	component.AddProperty("avg_review", "number", "Average reviews of the book", true)
	component.AddSimpleArrayProperty("image_urls", "string", "Book images url", true)
	component.AddProperty("authors", "string", "Comma-separated list of authors of the book", true)

	propertiesRefSchema := NewRefSchema("BookProperty")
	component.AddSchemaArrayProperty("properties", "Properties of the book", propertiesRefSchema, true)

	return component
}

func getBookPropertySchema() *Schema {
	component := NewSchema()

	component.AddProperty("id", "integer", "Generated id of the book property", true)
	component.AddProperty("isbn", "string", "Unique book's ISBN", true)
	component.AddProperty("book_id", "integer", "Generated id of the related book", true)
	component.AddProperty("format", "string", "Format of the related book", true)
	component.AddProperty("language", "string", "Language of the related book", true)
	component.AddProperty("price", "number", "Price of the related book in current format", true)
	component.AddProperty("publisher", "string", "Publisher of the related book", true)
	component.AddProperty("target_audience", "string", "Target audience of the related book", true)
	component.AddProperty("illustrator", "string", "Illustrator of the related book", true)
	component.AddProperty("illustrations", "string", "Illustrations of the related book", true)
	component.AddProperty("page_number", "integer", "Number of pages of the related book", true)
	component.AddProperty("available", "string", "Availability of the related book in current format", true)
	component.AddProperty("published_at", "string", "Publish date of the related book", true)

	return component
}
