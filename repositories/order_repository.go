package repository

import (
	"context"
	"fmt"
	customerror "marketplace-system/lib/customerrors"
	"marketplace-system/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type orderRepository repository

type OrderInterface interface {
	CreateOrder(ctx context.Context, tx *gorm.DB, order models.Order) (models.Order, error)
	CreateOrderDetail(ctx context.Context, tx *gorm.DB, orderDetail []*models.OrderDetail) ([]*models.OrderDetail, error)
	GetOrderById(ctx context.Context, id int) (order models.Order, err error)
	GetOrders(ctx context.Context, id int) (order []models.Order, err error)
	GetOrderByInvoice(ctx context.Context, invoiceNumber string) (order models.Order, err error)
}

func (c *orderRepository) CreateOrder(ctx context.Context, tx *gorm.DB, order models.Order) (models.Order, error) {
	err := tx.WithContext(ctx).Create(&order).Error
	// error record not found
	if err != nil {
		// error form another status
		logrus.Error(fmt.Sprintf("Err - create order - %s", err.Error()))
		return order, customerror.NewInternalError(err.Error())
	}
	return order, nil
}

func (c *orderRepository) CreateOrderDetail(ctx context.Context, tx *gorm.DB, orderDetail []*models.OrderDetail) ([]*models.OrderDetail, error) {
	err := tx.WithContext(ctx).Create(orderDetail).Error
	// error record not found
	if err != nil {
		// error form another status
		logrus.Error(fmt.Sprintf("Err - create order detail - %s", err.Error()))
		return orderDetail, customerror.NewInternalError(err.Error())
	}
	return orderDetail, nil
}

func (c *orderRepository) GetOrderById(ctx context.Context, id int) (order models.Order, err error) {
	err = c.Options.Postgres.Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Where("order_detail_status = ? ", "active")
	}).Where("order_id = ?", id).Where("order_status = ?", "pending").First(&order).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart list - %s", err.Error()))
		return order, err
	}

	return order, nil
}

func (c *orderRepository) GetOrders(ctx context.Context, id int) (order []models.Order, err error) {
	err = c.Options.Postgres.Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Where("order_detail_status = ? ", "active")
	}).Where("customer_id = ?", id).Find(&order).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get order list - %s", err.Error()))
		return order, err
	}

	return order, nil
}

func (c *orderRepository) GetOrderByInvoice(ctx context.Context, invoiceNumber string) (order models.Order, err error) {
	err = c.Options.Postgres.Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Where("order_detail_status = ? ", "active")
	}).Where("invoice_number = ?", invoiceNumber).First(&order).Error
	if err != nil {
		logrus.Error(fmt.Sprintf("Err - get cart list - %s", err.Error()))
		return order, err
	}

	return order, nil
}
