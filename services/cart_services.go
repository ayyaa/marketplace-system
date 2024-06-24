package services

import (
	"context"
	"encoding/json"
	"fmt"
	customerror "marketplace-system/lib/customerrors"
	"marketplace-system/models"
	"time"

	"github.com/go-redis/redis/v8"
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

	// Generate Redis key for the cart
	cartKey := fmt.Sprintf("cart:%d", addCart.CustomerID)
	tx := c.Options.Postgres.Begin()
	// Try to get the cart from Redis
	cartData, err := c.Options.Repository.Cart.GetCartListRedis(ctx, cartKey)
	var (
		cart      *models.Cart
		newDetail models.CartDetail
	)

	if err == redis.Nil {
		// Cart not found in Redis, fallback to DB
		cart, err = c.Options.Repository.Cart.GetOrCreateCart(ctx, tx, addCart)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// Cart found in Redis, unmarshal it
		if err = json.Unmarshal([]byte(cartData), &cart); err != nil {
			return err
		}
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
			return customerror.NewInternalError("Insufficient stock")
		}
		existingDetail.UpdatedAt = time.Now()
		existingDetail, err = c.Options.Repository.CartDetail.Save(ctx, tx, existingDetail)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Update the existing detail in the cart details
		for i, detail := range cart.Details {
			if detail.CartDetailID == existingDetail.CartDetailID {
				cart.Details[i] = existingDetail
				break
			}
		}
	} else {
		newDetail = models.CartDetail{
			CartID:           cart.CartID,
			ProductID:        addCart.ProductID,
			Quantity:         addCart.Quantity,
			CartDetailUUID:   uuid.New().String(),
			CartDetailStatus: "active",
		}
		cart.Details = append(cart.Details, newDetail)
		existingDetail, err = c.Options.Repository.CartDetail.Create(ctx, tx, newDetail)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	if err = c.Options.Repository.Cart.SetCartRedis(ctx, cartKey, *cart); err != nil {
		return err
	}

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

	// Generate Redis key for the cart
	cartKey := fmt.Sprintf("cart:%d", decreaseRequest.CustomerID)
	// Try to get the cart from Redis
	cartData, err := c.Options.Repository.Cart.GetCartListRedis(ctx, cartKey)
	var cart *models.Cart

	if err == redis.Nil {
		// Cart not found in Redis, fallback to DB
		cart, err = c.Options.Repository.Cart.GetCart(ctx, decreaseRequest.CustomerID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// Cart found in Redis, unmarshal it
		if err = json.Unmarshal([]byte(cartData), &cart); err != nil {
			return err
		}
	}

	decreaseRequest.CartID = cart.CartID
	// Check if the product is already in the cart
	var existingDetail models.CartDetail
	existingDetail, err = c.Options.Repository.CartDetail.CheckExistingCartDetail(ctx, decreaseRequest)
	if err != nil {
		return err
	}

	tx := c.Options.Postgres.Begin()
	// Update quantity if already in cart, otherwise return error
	if existingDetail.CartDetailID != 0 {
		existingDetail.Quantity -= decreaseRequest.Quantity
		existingDetail.UpdatedAt = time.Now()

		if existingDetail.Quantity <= 0 {
			// Check if product decrease until 0 and set status to deleted
			existingDetail.CartDetailStatus = "deleted_by_customer"
			existingDetail, err = c.Options.Repository.CartDetail.Save(ctx, tx, existingDetail)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Remove the detail from the cart
			for i, detail := range cart.Details {
				if detail.CartDetailID == existingDetail.CartDetailID {
					cart.Details = append(cart.Details[:i], cart.Details[i+1:]...)
					break
				}
			}

		} else {
			// update qty in cart detail
			existingDetail, err = c.Options.Repository.CartDetail.Save(ctx, tx, existingDetail)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Update the existing detail in the cart details
			for i, detail := range cart.Details {
				if detail.CartDetailID == existingDetail.CartDetailID {
					cart.Details[i] = existingDetail
					break
				}
			}
		}

		tx.Commit()

		if err = c.Options.Repository.Cart.SetCartRedis(ctx, cartKey, *cart); err != nil {
			return err
		}

		return nil
	}

	return customerror.NewInternalError("cannot decrease qty in cart")
}

func (c *cartServices) GetCartList(ctx context.Context, id int) (cart models.Cart, err error) {
	// Generate Redis key for the cart
	cartKey := fmt.Sprintf("cart:%d", id)

	// Try to get the cart from Redis
	cartData, err := c.Options.Repository.Cart.GetCartListRedis(ctx, cartKey)
	if err == nil {
		// Cart found in Redis, unmarshal it
		err = json.Unmarshal([]byte(cartData), &cart)
		if err != nil {
			return cart, err
		}
		return cart, nil
	} else if err != redis.Nil {
		// Redis error other than key not found
		return cart, err
	}

	// Cart not found in Redis, fallback to DB
	cart, err = c.Options.Repository.Cart.GetCartList(ctx, id)
	if err != nil {
		return cart, err
	}

	// Set cart data in Redis for future access
	err = c.Options.Repository.Cart.SetCartRedis(ctx, cartKey, cart)
	if err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cartServices) DeleteFromCart(ctx context.Context, deleteRequest models.ActionCart) (err error) {
	// Check if product exists
	var product models.Product
	product, err = c.Options.Repository.Product.GetProductBySlug(ctx, deleteRequest.ProductSlug)
	if err != nil {
		return err
	}

	deleteRequest.ProductID = product.ProductID

	// Generate Redis key for the cart
	cartKey := fmt.Sprintf("cart:%d", deleteRequest.CustomerID)
	// Try to get the cart from Redis
	cartData, err := c.Options.Repository.Cart.GetCartListRedis(ctx, cartKey)
	var cart *models.Cart

	if err == redis.Nil {
		// Cart not found in Redis, fallback to DB
		cart, err = c.Options.Repository.Cart.GetCart(ctx, deleteRequest.CustomerID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// Cart found in Redis, unmarshal it
		if err = json.Unmarshal([]byte(cartData), &cart); err != nil {
			return err
		}
	}

	tx := c.Options.Postgres.Begin()
	deleteRequest.CartID = cart.CartID
	// Check if the product is already in the cart detail
	var existingDetail models.CartDetail
	existingDetail, err = c.Options.Repository.CartDetail.CheckExistingCartDetail(ctx, deleteRequest)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update quantity if already in cart, otherwise create new detail
	if existingDetail.CartDetailID != 0 {
		existingDetail.UpdatedAt = time.Now()
		existingDetail.Quantity = 0
		existingDetail.CartDetailStatus = "deleted_by_customer"
		existingDetail, err = c.Options.Repository.CartDetail.Save(ctx, tx, existingDetail)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Remove the detail from the cart
		for i, detail := range cart.Details {
			if detail.CartDetailID == existingDetail.CartDetailID {
				cart.Details = append(cart.Details[:i], cart.Details[i+1:]...)
				break
			}
		}

		tx.Commit()

		// Otherwise, update the cart in Redis
		if err = c.Options.Repository.Cart.SetCartRedis(ctx, cartKey, *cart); err != nil {
			return err
		}

		return nil
	}

	return customerror.NewInternalError("product not found in cart")
}
