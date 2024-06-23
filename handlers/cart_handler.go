package handlers

import (
	"marketplace-system/lang"
	utils "marketplace-system/lib/helper"
	_validator "marketplace-system/lib/validator"
	"marketplace-system/models"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type cartHandlers handlers

type CartInterface interface {
	AddToCart(ctx echo.Context) error
	DecreaseFromCart(ctx echo.Context) error
	DeleteFromCart(ctx echo.Context) error
}

// Customer godoc
//
//	@Summary		Add to cart
//	@Description	Add to cart
//	@ID				add-to-cart
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CartRequest	true	"add to cart"
//	@Success		201		{object}	models.Response{data=[]models.Customer}
//	@Failure		400		{object}	models.BasicResponse{message=[]string}
//	@Failure		500		{object}	models.BasicResponse{messsage=[]string}
//	@Router			/cart/add [patch]
func (c *cartHandlers) AddToCart(ctx echo.Context) error {
	var req models.CartRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(err.Error(), http.StatusBadRequest))
	}

	// Validate struct
	err := _validator.Validate().Struct(req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			message := _validator.GetValidatorErrMsg(err.(validator.ValidationErrors))
			return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(message, http.StatusBadRequest))
		}
	}

	id := ctx.Get("id").(int)
	addToCart := models.ActionCart{
		ProductSlug: req.ProductSlug,
		Quantity:    req.Quantity,
		CustomerID:  id,
	}

	err = c.Options.Services.Cart.AddToCart(ctx.Request().Context(), addToCart)
	if err != nil {
		return utils.RespondWithError(ctx, http.StatusInternalServerError, utils.GetErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, nil)
}

// Cart godoc
//
//	@Summary		Decrease From cart
//	@Description	Decrease product from cart detail one by one or by qty
//	@ID				decrease-from-cart
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CartRequest	true	"add to cart"
//	@Success		201		{object}	models.Response{data=[]models.Customer}
//	@Failure		400		{object}	models.BasicResponse{message=[]string}
//	@Failure		500		{object}	models.BasicResponse{messsage=[]string}
//	@Router			/cart/decrease [patch]
func (c *cartHandlers) DecreaseFromCart(ctx echo.Context) error {
	var req models.CartRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(err.Error(), http.StatusBadRequest))
	}

	// Validate struct
	err := _validator.Validate().Struct(req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			message := _validator.GetValidatorErrMsg(err.(validator.ValidationErrors))
			return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(message, http.StatusBadRequest))
		}
	}

	id := ctx.Get("id").(int)
	addToCart := models.ActionCart{
		ProductSlug: req.ProductSlug,
		Quantity:    req.Quantity,
		CustomerID:  id,
	}

	err = c.Options.Services.Cart.DecreaseFromCart(ctx.Request().Context(), addToCart)
	if err != nil {
		return utils.RespondWithError(ctx, http.StatusInternalServerError, utils.GetErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, nil)
}

// Cart godoc
//
//	@Summary		Decrease From cart
//	@Description	Delete product from cart detail
//	@ID				decrease-from-cart
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.DeleteCartRequest	true	"add to cart"
//	@Success		201		{object}	models.Response{data=[]models.Customer}
//	@Failure		400		{object}	models.BasicResponse{message=[]string}
//	@Failure		500		{object}	models.BasicResponse{messsage=[]string}
//	@Router			/cart/delete [patch]
func (c *cartHandlers) DeleteFromCart(ctx echo.Context) error {
	var req models.DeleteCartRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(err.Error(), http.StatusBadRequest))
	}

	// Validate struct
	err := _validator.Validate().Struct(req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			message := _validator.GetValidatorErrMsg(err.(validator.ValidationErrors))
			return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(message, http.StatusBadRequest))
		}
	}

	id := ctx.Get("id").(int)
	removeProductCart := models.ActionCart{
		ProductSlug: req.ProductSlug,
		CustomerID:  id,
	}

	err = c.Options.Services.Cart.DeleteFromCart(ctx.Request().Context(), removeProductCart)
	if err != nil {
		return utils.RespondWithError(ctx, http.StatusInternalServerError, utils.GetErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, nil)
}
