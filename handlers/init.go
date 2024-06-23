package handlers

import (
	"marketplace-system/config"
	"marketplace-system/services"
)

type Main struct {
	Customer CustomerInterface
	Product  ProductInterface
	Cart     CartInterface
}

type handlers struct {
	Options Options
}

type Options struct {
	Config   *config.Config
	Services *services.Main
}

func Init(opts Options) *Main {
	hdrl := &handlers{opts}

	m := &Main{
		Customer: (*customerHandlers)(hdrl),
		Product:  (*productHandlers)(hdrl),
		Cart:     (*cartHandlers)(hdrl),
	}

	return m
}
