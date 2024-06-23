package repository

import (
	"marketplace-system/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Main struct {
	Customer   CustomerInterface
	Product    ProductInterface
	Cart       CartInterface
	CartDetail CartDetailInterface
	Order      OrderInterface
}

type repository struct {
	Options Options
}

type Options struct {
	Postgres *gorm.DB
	Redis    *redis.Client
	Config   *config.Config
}

func Init(opts Options) *Main {
	repo := &repository{opts}

	m := &Main{
		Customer:   (*customerRepository)(repo),
		Product:    (*productRepository)(repo),
		Cart:       (*cartRepository)(repo),
		CartDetail: (*cartDetailRepository)(repo),
		Order:      (*orderRepository)(repo),
	}

	return m
}
