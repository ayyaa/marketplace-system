package services

import (
	"marketplace-system/config"
	repository "marketplace-system/repositories"

	"gorm.io/gorm"
)

type Main struct {
	Customer CustomerInterface
	Product  ProductInterface
	Cart     CartInterface
	Checkout CheckoutInterface
}

type services struct {
	Options Options
}

type Options struct {
	Repository *repository.Main
	Config     *config.Config
	Postgres   *gorm.DB
}

func Init(opts Options) *Main {
	srvs := &services{opts}

	m := &Main{
		Customer: (*customerServices)(srvs),
		Product:  (*productServices)(srvs),
		Cart:     (*cartServices)(srvs),
		Checkout: (*checkoutServices)(srvs),
	}

	return m
}
