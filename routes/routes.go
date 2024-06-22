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

		register := v1.Group("/register")
		{
			register.POST("", handlers.Customer.InsertCustomer)
		}

		login := v1.Group("/login")
		{
			login.POST("", handlers.Customer.Login)
		}

		// masterClient := v1.Group("/master-client")
		{
			// masterClient.POST("", handlers.Client.Create)
			// masterClient.GET("", handlers.Client.Get)
			// masterClient.GET("/:uuid", handlers.Client.GetDetail)
			// masterClient.PATCH("/:uuid", handlers.Client.Update)
			// masterClient.DELETE("/:uuid", handlers.Client.Delete)

			// masterClient.PATCH("/:uuid/workhour", handlers.Workhour.Update)
			// masterClient.GET("/:uuid/workhour", handlers.Workhour.Get)
			// masterClient.GET("/:uuid/zone", handlers.Zone.GetZoneByClient)
		}

		// masterWorkhour := v1.Group("/master-workhour")
		{
			// masterWorkhour.POST("", handlers.Workhour.Create)
		}

	}
}
