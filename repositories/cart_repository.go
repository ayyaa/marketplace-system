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
	AddProductToCart(ctx context.Context, cartDetail models.CartDetail) error
	GetOrCreateCart(ctx context.Context, tx *gorm.DB, ddToCart models.ActionCart) (*models.Cart, error)
	GetCart(ctx context.Context, id int) (cart *models.Cart, err error)
	GetCartList(ctx context.Context, id int) (cart models.Cart, err error)
	Save(ctx context.Context, tx *gorm.DB, cart *models.Cart) (models.Cart, error)

	SetCartRedis(ctx context.Context, cacheKey string, cart models.Cart) error
	GetCartListRedis(ctx context.Context, cacheKey string) (cart models.Cart, err error)
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

func (c *cartRepository) GetCartListRedis(ctx context.Context, cacheKey string) (cart models.Cart, err error) {
	val, err := c.Options.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart list redis - %s", err.Error()))
		return cart, err
	}

	// Cache hit, return the cached value
	err = json.Unmarshal([]byte(val), &cart)
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart list redis - %s", "Failed to deserialize carts"))
		return cart, err
	}

	return cart, nil
}

func (c *cartRepository) SetCartRedis(ctx context.Context, cacheKey string, cart models.Cart) error {
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
