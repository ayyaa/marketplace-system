package routes

import (
	"marketplace-system/handlers"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRouter(e *echo.Echo, handlers *handlers.Main) {

	v1 := e.Group("/v1")
	{
		v1.GET("/swagger/*", echoSwagger.WrapHandler)

	}
}
