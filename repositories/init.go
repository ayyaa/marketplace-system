package repository

import (
	"marketplace-system/config"

	"gorm.io/gorm"
)

type Main struct {
	Customer CustomerInterface
}

type repository struct {
	Options Options
}

type Options struct {
	Postgres *gorm.DB
	Config   *config.Config
}

func Init(opts Options) *Main {
	repo := &repository{opts}

	m := &Main{
		Customer: (*customerRepository)(repo),
	}

	return m
}
