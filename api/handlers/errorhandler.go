package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// func ValidateStruct(c *gin.Context, obj interface{}) map[string]string {
// 	validate := validator.New()
// 	err := validate.Struct(obj)
// 	validationErrors := make(map[string]string)

// 	if err != nil {
// 		if errs, ok := err.(validator.ValidationErrors); ok {
// 			for _, fieldErr := range errs {
// 				validationErrors[fieldErr.Field()] = CustomErrorMessage(fieldErr)
// 			}
// 		}
// 	}

// 	return validationErrors
// }

func CustomErrorMessage[T any](fe validator.FieldError, obj T) string {
	jsonField := GetJSONFieldName(fe, obj)
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s harus diisi", jsonField)
	case "min":
		return fmt.Sprintf("%s minimal harus %s karakter", jsonField, fe.Param())
	// You can add more cases for other validation tags if needed.
	default:
		return fmt.Sprintf("%s tidak valid", jsonField)
	}
}
