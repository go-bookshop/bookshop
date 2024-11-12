package doc

func createCategoryOperation() *Operation {
	op := NewOperation("Create new category", "", []string{"Categories"})

	reqSchema := NewSchema()
	reqSchema.AddProperty("name", "string", "The name of the category", true)
	reqSchema.AddProperty("description", "string", "Description of the category", true)

	successSchema := NewSchema()
	successSchema.AddProperty("id", "integer", "Generated id of the category", true)
	successSchema.AddProperty("name", "string", "The name of the category", true)
	successSchema.AddProperty("description", "string", "Description of the category", true)

	requestErrorSchema := NewSchema()
	requestErrorSchema.AddProperty("errMsg", "string", "Specific error", true)

	validationErrorSchema := NewSchema()
	validationErrorSchema.AddProperty("errors", "object", "A map of errors with fields as keys", true)

	op.AddRequestSchema(reqSchema, "application/json", true)
	op.AddResponseSchema(successSchema, "201", "application/json", "Success")
	op.AddResponseSchema(requestErrorSchema, "400", "application/json", "Bad request")
	op.AddResponseSchema(validationErrorSchema, "422", "application/json", "Bad request")

	return op
}
