package handlers

import (
	"marketplace-system/config"
	"marketplace-system/services"
)

type Main struct{}

type handlers struct {
	Options Options
}

type Options struct {
	Config   *config.Config
	Services *services.Main
}

func Init(opts Options) *Main {
	// hdrl := &handlers{opts}

	m := &Main{}

	return m
}
