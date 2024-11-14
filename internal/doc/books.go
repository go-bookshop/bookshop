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

	refSchema := NewRefSchema("Book")
	successSchema.AddSchemaArrayProperty("data", "Paginated array of books", refSchema, true)

	requestErrorSchema := NewSchema()
	requestErrorSchema.AddProperty("errMsg", "string", "Specific error", true)

	op.AddRequestSchema(reqSchema, "application/json", true)
	op.AddResponseSchema(successSchema, "201", "application/json", "Success")
	op.AddResponseSchema(requestErrorSchema, "400", "application/json", "Bad request")

	return op
}

func getBookComponentSchema() *Schema {
	component := NewSchema()

	component.AddProperty("id", "integer", "Book unique identifies", true)
	component.AddProperty("title", "string", "Title of the book", true)
	component.AddProperty("avg_review", "number", "Average reviews of the book", true)
	component.AddSimpleArrayProperty("image_urls", "string", "Book images url", true)

	return component
}
