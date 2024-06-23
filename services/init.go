package services

import (
	"marketplace-system/config"
	repository "marketplace-system/repositories"
)

type Main struct {
	Customer CustomerInterface
	Product  ProductInterface
}

type services struct {
	Options Options
}

type Options struct {
	Repository *repository.Main
	Config     *config.Config
}

func Init(opts Options) *Main {
	srvs := &services{opts}

	m := &Main{
		Customer: (*customerServices)(srvs),
		Product:  (*productServices)(srvs),
	}

	return m
}
