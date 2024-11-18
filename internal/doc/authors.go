package doc

const AuthorsTag = "Authors"

func createAuthorOperation() *Operation {
	op := NewOperation("Create new author", "", []string{AuthorsTag})

	reqSchema := NewSchema()
	reqSchema.AddProperty("name", "string", "The name of the author", true)
	reqSchema.AddProperty("bio", "string", "A brief biography of the author", true)

	successSchema := NewSchema()
	successSchema.AddProperty("id", "integer", "Generated id of the author", false)
	successSchema.AddProperty("name", "string", "The name of the author", false)
	successSchema.AddProperty("bio", "string", "A brief biography of the author", false)

	requestErrorSchema := NewSchema()
	requestErrorSchema.AddProperty("errMsg", "string", "Specific error", false)

	validationErrorSchema := NewSchema()
	validationErrorSchema.AddProperty("errors", "object", "A map of errors with fields as keys", false)

	op.AddRequestSchema(reqSchema, "application/json", true)
	op.AddResponseSchema(successSchema, "201", "application/json", "Success")
	op.AddResponseSchema(requestErrorSchema, "400", "application/json", "Bad request")
	op.AddResponseSchema(validationErrorSchema, "422", "application/json", "Bad request")

	return op
}
