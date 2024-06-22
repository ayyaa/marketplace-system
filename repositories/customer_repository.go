package repository

import (
	"context"
	customerror "marketplace-system/lib/customerrors"
	"marketplace-system/models"

	"gorm.io/gorm"
)

type customerRepository repository

type CustomerInterface interface {
	InsertCustomer(ctx context.Context, customer models.Customer) (models.Customer, error)
	GetCustomerLoginByPhone(ctx context.Context, customerReq models.LoginRequest) (models.Customer, error)
}

func (c *customerRepository) InsertCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
	err := c.Options.Postgres.WithContext(ctx).Create(&customer).Error
	// error record not found
	if err != nil {
		// error form another status
		return customer, customerror.NewInternalError(err.Error())
	}
	return customer, nil
}

func (c *customerRepository) GetCustomerLoginByPhone(ctx context.Context, customerReq models.LoginRequest) (models.Customer, error) {
	customer := models.Customer{}
	err := c.Options.Postgres.WithContext(ctx).Where("phone = ?", customerReq.Phone).First(&customer).Error
	// error record not found
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return customer, customerror.NewBadRequestErrorf(err.Error())
		}
		// error form another status
		return customer, customerror.NewInternalError(err.Error())
	}
	return customer, nil
}
