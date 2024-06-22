package services

import (
	"context"
	"marketplace-system/lang"
	customerror "marketplace-system/lib/customerrors"
	"marketplace-system/lib/helper"
	"marketplace-system/middleware"
	"marketplace-system/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type customerServices services

type CustomerInterface interface {
	InsertCustomer(ctx context.Context, customer models.CustomerRequest) (models.Customer, error)
	GetCustomerLoginByPhone(ctx context.Context, customerReq models.LoginRequest) (models.Customer, error)
	Login(ctx context.Context, customerReq models.LoginRequest) (string, error)
}

func (c *customerServices) InsertCustomer(ctx context.Context, customerReq models.CustomerRequest) (models.Customer, error) {
	var customer models.Customer

	// Generate UUID
	customer.CustomerUUID = uuid.New().String()

	// Hash the password
	hashedPassword, err := helper.HashPassword(customerReq.Password)
	if err != nil {
		return customer, customerror.NewInternalErrorf("failed to hash password")
	}
	customer.Password = hashedPassword

	// Set customer status
	customer.CustomerStatus = "active"
	customer.Address = customerReq.Address
	customer.FullName = customerReq.FullName
	customer.Email = customerReq.Email
	customer.Phone = customerReq.Phone

	// Insert customer into repository
	return c.Options.Repository.Customer.InsertCustomer(ctx, customer)

}

func (c *customerServices) GetCustomerLoginByPhone(ctx context.Context, customerReq models.LoginRequest) (models.Customer, error) {
	return c.Options.Repository.Customer.GetCustomerLoginByPhone(ctx, customerReq)

}

func (c *customerServices) Login(ctx context.Context, customerReq models.LoginRequest) (string, error) {
	// validate if any user by phone
	customer, err := c.GetCustomerLoginByPhone(ctx, customerReq)
	if err != nil {
		return "", err
	}

	// if no sql rows, user not found
	if customer.CustomerID == 0 {
		logrus.Error(err)
		return "", lang.ErrDataNotFound

	}

	// compare password from request and database
	match := helper.CheckPasswordHash(customerReq.Password, customer.Password)
	if match {
		// generate token auth
		token, err := middleware.GenerateToken(customer, *c.Options.Config)
		if err != nil {
			return "", err
		}

		// return succes and user token
		return token, nil
	}

	return "", customerror.NewInternalErrorf(lang.ErrorUnsuccessfulLogin)

}
