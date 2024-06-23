package repository

import (
	"context"
	"fmt"
	"marketplace-system/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type cartDetailRepository repository

type CartDetailInterface interface {
	CheckExistingCartDetail(ctx context.Context, addToCart models.ActionCart) (models.CartDetail, error)
	Save(ctx context.Context, tx *gorm.DB, existingDetail models.CartDetail) (models.CartDetail, error)
	Create(ctx context.Context, tx *gorm.DB, newDetail models.CartDetail) (models.CartDetail, error)
}

func (c *cartDetailRepository) CheckExistingCartDetail(ctx context.Context, addToCart models.ActionCart) (models.CartDetail, error) {
	var existingDetail models.CartDetail
	err := c.Options.Postgres.Where("cart_id = ? AND product_id = ? AND cart_detail_status = ?", addToCart.CartID, addToCart.ProductID, "active").First(&existingDetail).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logrus.Error(fmt.Sprintf("Err - get cart detail - %s", err.Error()))
			return existingDetail, err
		}
	}

	return existingDetail, nil
}

func (c *cartDetailRepository) Save(ctx context.Context, tx *gorm.DB, existingDetail models.CartDetail) (models.CartDetail, error) {
	err := tx.Save(&existingDetail).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart detail - %s", "Failed to update cart detail"))
		return existingDetail, err
	}

	return existingDetail, nil
}

func (c *cartDetailRepository) Create(ctx context.Context, tx *gorm.DB, newDetail models.CartDetail) (models.CartDetail, error) {
	err := tx.Create(&newDetail).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart detail - %s", "Failed to save cart detail"))
		return newDetail, err
	}

	return newDetail, nil
}
