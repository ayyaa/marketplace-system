package services

import (
	"context"
	"fmt"
	"marketplace-system/models"

	"github.com/go-redis/redis/v8"
)

type productServices services

type ProductInterface interface {
	FindProductsByCategory(ctx context.Context, categorySlug string) (products []models.Product, err error)
}

func (p *productServices) FindProductsByCategory(ctx context.Context, categorySlug string) (products []models.Product, err error) {

	products, err = p.Options.Repository.Product.GetProductsByCategoryRedis(ctx, categorySlug)
	if err == redis.Nil {
		// Cache miss, query the database
		category, err := p.Options.Repository.Product.FindProductsByCategory(ctx, categorySlug)
		if err != nil {
			return nil, err
		}
		cacheKey := fmt.Sprintf("products:category:%s", categorySlug)
		err = p.Options.Repository.Product.SetProductRedis(ctx, cacheKey, category.Products)
		if err != nil {
			return nil, err
		}

		return category.Products, nil
	} else if err != nil {
		return nil, err
	}

	return products, nil
}
