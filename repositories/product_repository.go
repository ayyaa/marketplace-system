package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace-system/models"
	"time"

	"github.com/sirupsen/logrus"
)

type productRepository repository

type ProductInterface interface {
	FindProductsByCategory(ctx context.Context, categorySlug string) (models.Category, error)
	GetProductsByCategoryRedis(ctx context.Context, categorySlug string) (products []models.Product, err error)
	GetProductBySlug(ctx context.Context, productSlug string) (models.Product, error)
	GetProductRedis(ctx context.Context, productID int) (models.Product, error)
	SetProductRedis(ctx context.Context, cacheKey string, products []models.Product) error
	GetProductById(ctx context.Context, id int) (models.Product, error)
}

func (c *productRepository) FindProductsByCategory(ctx context.Context, categorySlug string) (models.Category, error) {
	// Query products by category ID
	var category models.Category
	if err := c.Options.Postgres.Preload("Products").Where("category_slug = ?", categorySlug).First(&category).Error; err != nil {
		logrus.Error(fmt.Sprintf("Err - get product by category - %s", err.Error()))
		return category, err
	}

	return category, nil
}

func (c *productRepository) GetProductsByCategoryRedis(ctx context.Context, categorySlug string) (products []models.Product, err error) {
	// Check Redis cache
	cacheKey := fmt.Sprintf("products:category:%s", categorySlug)
	val, err := c.Options.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get product by category from redis - %s", err.Error()))
		return nil, err
	}

	// Cache hit, return the cached value
	err = json.Unmarshal([]byte(val), &products)
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get product by category from redis - %s", "Failed to deserialize products"))
		return nil, err
	}

	return products, nil
}

func (c *productRepository) GetProductBySlug(ctx context.Context, productSlug string) (models.Product, error) {
	var product models.Product
	if err := c.Options.Postgres.WithContext(ctx).Where("product_slug = ?", productSlug).First(&product).Error; err != nil {
		logrus.Error(fmt.Sprintf("Err - get product by slug - %s", err.Error()))
		return product, err
	}

	return product, nil
}

func (c *productRepository) GetProductById(ctx context.Context, id int) (models.Product, error) {
	var product models.Product
	if err := c.Options.Postgres.WithContext(ctx).Where("product_id = ?", id).First(&product).Error; err != nil {
		logrus.Error(fmt.Sprintf("Err - get product by id - %s", err.Error()))
		return product, err
	}

	return product, nil
}

func (c *productRepository) GetProductRedis(ctx context.Context, productID int) (models.Product, error) {
	var product models.Product
	// Check Redis cache
	cacheKey := fmt.Sprintf("products:%d", productID)
	val, err := c.Options.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get product by slug from redis - %s", err.Error()))
		return product, err
	}

	// Cache hit, return the cached value
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get product by slug from redis - %s", "Failed to deserialize products"))
		return product, err
	}

	return product, nil
}

func (c *productRepository) SetProductRedis(ctx context.Context, cacheKey string, products []models.Product) error {
	// Serialize products to JSON
	productsJSON, err := json.Marshal(products)
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - set product by slug from redis - %s", "Failed to serialize products"))
		return err
	}

	// Store the result in Redis
	err = c.Options.Redis.Set(ctx, cacheKey, productsJSON, 72*time.Hour).Err()
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - set product by slug from redis - %s", "Failed to store in Redis"))
		return err
	}

	return nil
}
