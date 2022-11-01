package utils

import (
	"github.com/go-playground/validator"
	"github.com/sabilhaq/devcode-todo-go/models"
)

var validate = validator.New()

func ValidateStruct(data interface{}) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			var errorField string
			if err.Field() == "Title" {
				errorField = "title"
			} else if err.Field() == "ActivityGroupID" {
				errorField = "activity_group_id"
			}

			element.FailedField = errorField
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
