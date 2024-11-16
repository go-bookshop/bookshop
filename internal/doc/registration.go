package doc

func registerUserOperation() *Operation {
	op := NewOperation("Register new user", "", []string{"Registration"})

	reqSchema := NewSchema()
	reqSchema.AddProperty("firstName", "string", "The first name of the user", true)
	reqSchema.AddProperty("lastName", "string", "The last name of the user", true)
	reqSchema.AddProperty("email", "string", "The email of the user", true)
	reqSchema.AddProperty("password", "string", "The Password for the user", true)

	successSchema := NewSchema()
	successSchema.AddProperty("id", "integer", "Generated id of the user", true)
	successSchema.AddProperty("firstName", "string", "The first name of the user", true)
	successSchema.AddProperty("lastName", "string", "The last name of the user", true)
	successSchema.AddProperty("email", "string", "The email of the user", true)
	successSchema.AddProperty("activated", "boolean", "Indicates whether the user has been activated via email verification", true)

	requestErrorSchema := NewSchema()
	requestErrorSchema.AddProperty("errMsg", "string", "Specific error", true)

	validationErrorSchema := NewSchema()
	validationErrorSchema.AddProperty("errors", "object", "A map of errors with fields as keys", true)

	op.AddRequestSchema(reqSchema, "application/json", true)
	op.AddResponseSchema(successSchema, "202", "application/json", "Success, sending activation email in the background")
	op.AddResponseSchema(requestErrorSchema, "400", "application/json", "Bad request")
	op.AddResponseSchema(validationErrorSchema, "422", "application/json", "Bad request")

	return op
}

func activateUserOperation() *Operation {
	op := NewOperation("Activate user", "Provide token to activate user", []string{"Registration"})

	reqSchema := NewSchema()
	reqSchema.AddProperty("token", "string", "Token for user activation", true)

	successSchema := NewSchema()
	successSchema.AddProperty("id", "integer", "Generated id of the user", true)
	successSchema.AddProperty("firstName", "string", "The first name of the user", true)
	successSchema.AddProperty("lastName", "string", "The last name of the user", true)
	successSchema.AddProperty("email", "string", "The email of the user", true)
	successSchema.AddProperty("activated", "boolean", "Indicates whether the user has been activated via email verification", true)

	requestErrorSchema := NewSchema()
	requestErrorSchema.AddProperty("errMsg", "string", "Specific error", true)

	validationErrorSchema := NewSchema()
	validationErrorSchema.AddProperty("errors", "object", "A map of errors with fields as keys", true)

	op.AddRequestSchema(reqSchema, "application/json", true)
	op.AddResponseSchema(successSchema, "200", "application/json", "Success")
	op.AddResponseSchema(requestErrorSchema, "400", "application/json", "Bad request")
	op.AddResponseSchema(requestErrorSchema, "409", "application/json", "Race condition on activating user")
	op.AddResponseSchema(validationErrorSchema, "422", "application/json", "Bad request")

	return op
}

func resendTokenOperation() *Operation {
	op := NewOperation("Resend token", "Resend token to activate user", []string{"Registration"})

	reqSchema := NewSchema()
	reqSchema.AddProperty("email", "string", "Registered but not activated email", true)

	requestErrorSchema := NewSchema()
	requestErrorSchema.AddProperty("errMsg", "string", "Specific error", true)

	validationErrorSchema := NewSchema()
	validationErrorSchema.AddProperty("errors", "object", "A map of errors with fields as keys", true)

	op.AddRequestSchema(reqSchema, "application/json", true)
	op.AddResponseSchema(NoBodySchema, "202", "", "Success, sending activation email in the background")
	op.AddResponseSchema(requestErrorSchema, "400", "application/json", "Bad request")
	op.AddResponseSchema(validationErrorSchema, "422", "application/json", "Bad request")

	return op
}
