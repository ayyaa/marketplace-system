package services

import (
	"marketplace-system/config"
	repository "marketplace-system/repositories"
)

type Main struct {
}

type services struct {
	Options Options
}

type Options struct {
	Repository *repository.Main
	Config     *config.Config
}

func Init(opts Options) *Main {
	// ucs := &services{opts}

	m := &Main{}

	return m
}
