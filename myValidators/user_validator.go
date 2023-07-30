package myValidators

import "github.com/go-playground/validator/v10"

var FieldErrorMessages = map[string]map[string]string{
	"FirstName": {
		"required": "First name is required",
		"min":      "First name should have a minimum of 3 characters",
		"max":      "First name should have a maximum of 15 characters",
	},
	"LastName": {
		"required": "Last name is required",
		"min":      "Last name should have a minimum of 3 characters",
		"max":      "Last name should have a maximum of 15 characters",
	},
	"Email": {
		"required": "Email is required",
		"email":    "Invalid email address",
	},
}

func ExtractUserValidationErrors(err error) map[string]string {
	validationErrors := err.(validator.ValidationErrors)
	errorMessages := make(map[string]string)

	// Custom error messages for specific fields and validation rules
	for _, err := range validationErrors {
		fieldName := err.Field()
		tagName := err.Tag()
		if customMessage, ok := FieldErrorMessages[fieldName][tagName]; ok {
			errorMessages[fieldName] = customMessage
		}
	}

	return errorMessages
}
