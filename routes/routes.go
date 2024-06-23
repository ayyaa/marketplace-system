package routes

import (
	"marketplace-system/config"
	"marketplace-system/handlers"
	"marketplace-system/middleware"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRouter(e *echo.Echo, handlers *handlers.Main, cfg *config.Config) {

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

		categories := v1.Group("/categories")
		{
			categories.GET("/:slug/products", handlers.Product.FindProductsByCategory)
		}

		cart := v1.Group("/cart", middleware.JWTMiddleware(cfg.SecretKeyJWT))
		{
			cart.PATCH("/add", handlers.Cart.AddToCart)
			cart.PATCH("/decrease", handlers.Cart.DecreaseFromCart)
			cart.PATCH("/delete", handlers.Cart.DeleteFromCart)
			cart.GET("", handlers.Cart.CartDetailList)
		}

		checkout := v1.Group("/checkout", middleware.JWTMiddleware(cfg.SecretKeyJWT))
		{
			checkout.POST("", handlers.Checkout.Checkout)
		}

	}
}
