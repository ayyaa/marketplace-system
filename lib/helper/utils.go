package helper

import (
	"fmt"
	"marketplace-system/lang"
	"marketplace-system/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ResponseSuccess for sending response with data
func ResponseSuccess(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, models.ResponseSuccess{
		Message: lang.SuccessMsg,
		Success: true,
		Data:    data,
	})
}

// RespondWithError for sending response with status not 200
func RespondWithError(c echo.Context, status int, data interface{}) error {
	c.Set("responseBody", data)
	return c.JSON(status, data)
}

// GetErrorResponse is a function for build errorMessage
func GetErrorResponse(message string, statusCode int) *models.ApplicationError {
	return &models.ApplicationError{
		Message:    fmt.Sprintf(message),
		StatusCode: statusCode,
		Success:    false,
	}
}
