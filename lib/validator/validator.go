package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	UnmappedrrMsg map[string]string = map[string]string{
		"required": "Value is required",
		"max":      "Value is more than max character",
		"min":      "value is less thank min character",
		"email":    "value is invalid",
	}
)

var _validator *validator.Validate

func Validate() *validator.Validate {
	return _validator
}

func InitValidator(extendedFuncs ...func(*validator.Validate)) {
	_validator = validator.New()

	for _, extFunc := range extendedFuncs {
		extFunc(_validator)
	}
}

func ValidateErrToMapString(valErr validator.ValidationErrors) map[string]string {
	result := make(map[string]string, len(valErr))
	for _, err := range valErr {
		result[err.Field()] = fmt.Sprintf("validation: %v", err.ActualTag())
	}
	return result
}

func GetValidatorErrMsg(e validator.ValidationErrors) string {
	var strArr []string
	for _, v := range e {
		newstr := fmt.Sprintf("%s %s", v.Field(), UnmappedrrMsg[v.Tag()])
		strArr = append(strArr, newstr)
	}

	return strings.Join(strArr, ", ")
}
