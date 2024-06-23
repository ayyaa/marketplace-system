package services

import (
	"context"
	"fmt"
	customerror "marketplace-system/lib/customerrors"
	"marketplace-system/models"
	"time"

	"github.com/google/uuid"
)

type checkoutServices services

type CheckoutInterface interface {
	Checkout(ctx context.Context, checkout models.Checkout) (order models.Order, err error)
}

func (c *checkoutServices) Checkout(ctx context.Context, checkout models.Checkout) (order models.Order, err error) {
	var (
		orderDetails []*models.OrderDetail
	)
	cart, err := c.Options.Repository.Cart.GetCartList(ctx, checkout.CustomerID)
	if err != nil {
		return order, customerror.NewNotFoundError("Cart Not Found")
	}

	if len(cart.Details) < 0 {
		return order, customerror.NewNotFoundError("Cart Details Not Found")
	}

	grandTotal := 0.0
	for _, detail := range cart.Details {
		product, err := c.Options.Repository.Product.GetProductById(ctx, detail.ProductID)
		if err != nil {
			return order, err
		}
		orderDetail := models.OrderDetail{
			OrderDetailUUID:   uuid.New().String(),
			ProductID:         detail.ProductID,
			Quantity:          detail.Quantity,
			Price:             product.Price,
			FinalPrice:        float64(detail.Quantity) * product.Price,
			OrderDetailStatus: "active",
		}
		orderDetails = append(orderDetails, &orderDetail)
		grandTotal += float64(detail.Quantity) * product.Price
	}

	tx := c.Options.Postgres.Begin()

	order = models.Order{
		OrderUUID:          uuid.New().String(),
		InvoiceNumber:      fmt.Sprintf("INV-%d", time.Now().Unix()),
		CustomerID:         cart.CustomerID,
		CartID:             cart.CartID,
		OrderPaymentType:   checkout.PaymentType,
		OrderPaymentStatus: "unpaid",
		OrderStatus:        "pending",
		GrandTotal:         grandTotal,
		ExpiredAt:          time.Now().Add(24 * time.Hour),
	}

	order, err = c.Options.Repository.Order.CreateOrder(ctx, tx, order)
	if err != nil {
		tx.Rollback()
		return order, err
	}

	for _, detail := range orderDetails {
		detail.OrderID = order.OrderID
	}

	_, err = c.Options.Repository.Order.CreateOrderDetail(ctx, tx, orderDetails)
	if err != nil {
		tx.Rollback()
		return order, err
	}

	cart.CartStatus = "converted"
	cart, err = c.Options.Repository.Cart.Save(ctx, tx, &cart)
	if err != nil {
		tx.Rollback()
		return order, err
	}

	tx.Commit()

	order, err = c.Options.Repository.Order.GetOrderById(ctx, order.OrderID)
	if err != nil {
		return order, err
	}

	return order, nil
}
