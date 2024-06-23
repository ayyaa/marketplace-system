package repository

import (
	"context"
	"fmt"
	"marketplace-system/models"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type cartRepository repository

type CartInterface interface {
	AddProductToCart(ctx context.Context, cartDetail models.CartDetail) error
	GetOrCreateCart(ctx context.Context, tx *gorm.DB, ddToCart models.ActionCart) (*models.Cart, error)
	GetCart(ctx context.Context, addToCart models.ActionCart) (cart *models.Cart, err error)
}

func (c *cartRepository) AddProductToCart(ctx context.Context, cartDetail models.CartDetail) error {
	return c.Options.Postgres.WithContext(ctx).Create(&cartDetail).Error
}

func (c *cartRepository) GetOrCreateCart(ctx context.Context, tx *gorm.DB, addToCart models.ActionCart) (*models.Cart, error) {
	// Create or update cart
	cart := models.Cart{
		CartUUID:   uuid.New().String(),
		CartStatus: "active",
		CustomerID: addToCart.CustomerID,
		UpdatedAt:  time.Now(),
	}

	err := tx.Where("customer_id = ?", addToCart.CustomerID).Where("cart_status = ?", "active").FirstOrCreate(&cart).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart - %s", err.Error()))
		return nil, err
	}

	return &cart, nil
}

func (c *cartRepository) GetCart(ctx context.Context, addToCart models.ActionCart) (cart *models.Cart, err error) {
	err = c.Options.Postgres.Where("customer_id = ?", addToCart.CustomerID).Where("cart_status = ?", "active").First(&cart).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart - %s", err.Error()))
		return nil, err
	}

	return cart, nil
}
