package services

import (
	"context"
	"marketplace-system/models"
)

type productServices services

type ProductInterface interface {
	FindProductsByCategory(ctx context.Context, categorySlug string) (products []models.Product, err error)
}

func (p *productServices) FindProductsByCategory(ctx context.Context, categorySlug string) (products []models.Product, err error) {

	// Cache miss, query the database
	category, err := p.Options.Repository.Product.FindProductsByCategory(ctx, categorySlug)
	if err != nil {
		return nil, err
	}

	return category.Products, nil

}
