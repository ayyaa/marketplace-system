package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace-system/models"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type cartRepository repository

type CartInterface interface {
	GetOrCreateCart(ctx context.Context, tx *gorm.DB, ddToCart models.ActionCart) (*models.Cart, error)
	GetCart(ctx context.Context, id int) (cart *models.Cart, err error)
	GetCartList(ctx context.Context, id int) (cart models.Cart, err error)
	Save(ctx context.Context, tx *gorm.DB, cart *models.Cart) (models.Cart, error)

	SetCartRedis(ctx context.Context, cacheKey string, cart models.Cart) error
	GetCartListRedis(ctx context.Context, cacheKey string) (cart string, err error)
	DeleteCartRedis(ctx context.Context, cacheKey string) error
}

func (c *cartRepository) GetOrCreateCart(ctx context.Context, tx *gorm.DB, addToCart models.ActionCart) (*models.Cart, error) {
	// Create or update cart
	cart := models.Cart{
		CartUUID:   uuid.New().String(),
		CartStatus: "active",
		CustomerID: addToCart.CustomerID,
		UpdatedAt:  time.Now(),
	}

	err := tx.Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Where("cart_detail_status = ? ", "active")
	}).Where("customer_id = ?", addToCart.CustomerID).Where("cart_status = ?", "active").FirstOrCreate(&cart).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart - %s", err.Error()))
		return nil, err
	}

	return &cart, nil
}

func (c *cartRepository) GetCart(ctx context.Context, id int) (cart *models.Cart, err error) {
	err = c.Options.Postgres.Where("customer_id = ?", id).Where("cart_status = ?", "active").First(&cart).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart - %s", err.Error()))
		return nil, err
	}

	return cart, nil
}

func (c *cartRepository) GetCartList(ctx context.Context, id int) (cart models.Cart, err error) {
	err = c.Options.Postgres.Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Where("cart_detail_status = ? ", "active")
	}).Where("customer_id = ?", id).Where("cart_status = ?", "active").First(&cart).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart list - %s", err.Error()))
		return cart, err
	}

	return cart, nil
}

func (c *cartRepository) Save(ctx context.Context, tx *gorm.DB, cart *models.Cart) (models.Cart, error) {
	err := tx.Save(&cart).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - updated cart - %s", "Failed to update cart"))
		return *cart, err
	}

	return *cart, nil
}

func (c *cartRepository) GetCartListRedis(ctx context.Context, cacheKey string) (cart string, err error) {
	val, err := c.Options.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart list redis - %s", err.Error()))
		return "", err
	}

	return val, nil
}

func (c *cartRepository) SetCartRedis(ctx context.Context, cacheKey string, cart models.Cart) error {
	if err := c.DeleteCartRedis(ctx, cacheKey); err != nil {
		logrus.Error(fmt.Sprintf("Err - set cart redis - %s", "Failed to remove cart redis"))
		return err
	}
	// Serialize products to JSON
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - set cart redis - %s", "Failed to serialize cart"))
		return err
	}

	// Store the result in Redis
	err = c.Options.Redis.Set(ctx, cacheKey, cartJSON, 72*time.Hour).Err()
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - set cart redis - %s", "Failed to store in Redis"))
		return err
	}

	return nil
}

func (c *cartRepository) DeleteCartRedis(ctx context.Context, cacheKey string) error {
	if err := c.Options.Redis.Del(ctx, cacheKey).Err(); err != nil {
		logrus.Error(fmt.Sprintf("Err - delete cart redis - %s", "Failed to remove cart redis"))
		return err
	}

	return nil
}
