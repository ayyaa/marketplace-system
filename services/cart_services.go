package services

import (
	"context"
	customerror "marketplace-system/lib/customerrors"
	"marketplace-system/models"
	"time"

	"github.com/google/uuid"
)

type cartServices services

type CartInterface interface {
	AddToCart(ctx context.Context, addCart models.ActionCart) (err error)
	DecreaseFromCart(ctx context.Context, addCart models.ActionCart) (err error)
	DeleteFromCart(ctx context.Context, addCart models.ActionCart) (err error)
	GetCartList(ctx context.Context, id int) (models.Cart, error)
}

func (c *cartServices) AddToCart(ctx context.Context, addCart models.ActionCart) (err error) {
	// Check if product exists
	var product models.Product
	product, err = c.Options.Repository.Product.GetProductBySlug(ctx, addCart.ProductSlug)
	if err != nil {
		return err
	}

	// Check if quantity is available in stock
	if product.StockQuantity < addCart.Quantity {
		return customerror.NewInternalError("Insufficient stock")
	}

	addCart.ProductID = product.ProductID

	tx := c.Options.Postgres.Begin()
	// Create or update cart
	cart, err := c.Options.Repository.Cart.GetOrCreateCart(ctx, tx, addCart)
	if err != nil {
		tx.Rollback()
		return err
	}

	addCart.CartID = cart.CartID
	// Check if the product is already in the cart
	var existingDetail models.CartDetail
	existingDetail, err = c.Options.Repository.CartDetail.CheckExistingCartDetail(ctx, addCart)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update quantity if already in cart, otherwise create new detail
	if existingDetail.CartDetailID != 0 {
		existingDetail.Quantity += addCart.Quantity
		// Check if quantity is available in stock
		if product.StockQuantity < existingDetail.Quantity {
			tx.Rollback()
			return customerror.NewInternalError("Insufficient stock")
		}
		existingDetail.UpdatedAt = time.Now()
		existingDetail, err = c.Options.Repository.CartDetail.Save(ctx, tx, existingDetail)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		newDetail := models.CartDetail{
			CartID:           cart.CartID,
			ProductID:        addCart.ProductID,
			Quantity:         addCart.Quantity,
			CartDetailUUID:   uuid.New().String(),
			CartDetailStatus: "active",
		}
		existingDetail, err = c.Options.Repository.CartDetail.Create(ctx, tx, newDetail)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (c *cartServices) DecreaseFromCart(ctx context.Context, decreaseRequest models.ActionCart) (err error) {
	// Check if product exists
	var product models.Product
	product, err = c.Options.Repository.Product.GetProductBySlug(ctx, decreaseRequest.ProductSlug)
	if err != nil {
		return err
	}

	decreaseRequest.ProductID = product.ProductID
	// Create or update cart
	cart, err := c.Options.Repository.Cart.GetCart(ctx, decreaseRequest)
	if err != nil {
		return err
	}

	decreaseRequest.CartID = cart.CartID
	// Check if the product is already in the cart
	var existingDetail models.CartDetail
	existingDetail, err = c.Options.Repository.CartDetail.CheckExistingCartDetail(ctx, decreaseRequest)
	if err != nil {
		return err
	}

	// Update quantity if already in cart, otherwise create new detail
	if existingDetail.CartDetailID != 0 {
		existingDetail.Quantity -= decreaseRequest.Quantity
		existingDetail.UpdatedAt = time.Now()

		if existingDetail.Quantity <= 0 {
			// Check if product decrease until 0 and set status to deleted
			existingDetail.CartDetailStatus = "deleted_by_customer"
			existingDetail, err = c.Options.Repository.CartDetail.Save(ctx, c.Options.Postgres, existingDetail)
			if err != nil {
				return err
			}
			return nil
		}

		// update qty in cart detail
		existingDetail, err = c.Options.Repository.CartDetail.Save(ctx, c.Options.Postgres, existingDetail)
		if err != nil {
			return err
		}

		return nil
	}

	return customerror.NewInternalError("cannot decrease qty in cart")
}

func (c *cartServices) DeleteFromCart(ctx context.Context, deleteRequest models.ActionCart) (err error) {
	// Check if product exists
	var product models.Product
	product, err = c.Options.Repository.Product.GetProductBySlug(ctx, deleteRequest.ProductSlug)
	if err != nil {
		return err
	}

	deleteRequest.ProductID = product.ProductID
	cart, err := c.Options.Repository.Cart.GetCart(ctx, deleteRequest)
	if err != nil {
		return err
	}

	deleteRequest.CartID = cart.CartID
	// Check if the product is already in the cart detail
	var existingDetail models.CartDetail
	existingDetail, err = c.Options.Repository.CartDetail.CheckExistingCartDetail(ctx, deleteRequest)
	if err != nil {
		return err
	}

	// Update quantity if already in cart, otherwise create new detail
	if existingDetail.CartDetailID != 0 {
		existingDetail.UpdatedAt = time.Now()
		existingDetail.Quantity = 0
		existingDetail.CartDetailStatus = "deleted_by_customer"
		existingDetail, err = c.Options.Repository.CartDetail.Save(ctx, c.Options.Postgres, existingDetail)
		if err != nil {
			return err
		}

		return nil
	}

	return customerror.NewInternalError("product not found in cart")
}

func (c *cartServices) GetCartList(ctx context.Context, id int) (cart models.Cart, err error) {
	cart, err = c.Options.Repository.Cart.GetCartList(ctx, id)
	if err != nil {
		return cart, err
	}

	return cart, nil
}
