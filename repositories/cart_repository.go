package repository

import (
	"context"
	"marketplace-system/models"
)

type cartRepository repository

type CartInterface interface {
	AddProductToCart(ctx context.Context, cartDetail models.CartDetail) error
}

func (c *cartRepository) AddProductToCart(ctx context.Context, cartDetail models.CartDetail) error {
	return c.Options.Postgres.WithContext(ctx).Create(&cartDetail).Error
}
