package val

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/zvash/accounting-system/internal/util"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		Validator *validator.Validate
	}
)

// Validate This is the validator instance
// for more information see: https://github.com/go-playground/validator
var Validate = validator.New()

func NewValidator() *XValidator {
	xValidator := &XValidator{
		Validator: Validate,
	}
	xValidator.RegisterCustomValidators()
	return xValidator
}

func (v XValidator) RegisterCustomValidators() {
	_ = Validate.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		if currency, ok := fl.Field().Interface().(string); ok {
			return util.IsSupportedCurrency(currency)
		}
		return false
	})
}

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	errs := Validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func (v XValidator) MakeErrorBag(errs []ErrorResponse) []string {
	errorMessages := make([]string, 0)

	for _, err := range errs {
		errorMessages = append(errorMessages, fmt.Sprintf(
			"[%s]: '%v' | Needs to implement '%s'",
			err.FailedField,
			err.Value,
			err.Tag,
		))
	}
	return errorMessages
}
