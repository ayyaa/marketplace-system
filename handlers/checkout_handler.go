package handlers

import (
	"marketplace-system/lang"
	customerror "marketplace-system/lib/customerrors"
	utils "marketplace-system/lib/helper"
	_validator "marketplace-system/lib/validator"
	"marketplace-system/models"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type checkoutHandlers handlers

type CheckoutInterface interface {
	Checkout(ctx echo.Context) error
}

// Customer godoc
//
//	@Summary		Checkout order
//	@Description	Checkout order
//	@ID				checkout-order
//	@Tags			Checkout
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CustomerRequest	true	"Checkout Order"
//	@Success		201		{object}	models.Response{data=[]models.Customer}
//	@Failure		400		{object}	models.BasicResponse{message=[]string}
//	@Failure		500		{object}	models.BasicResponse{messsage=[]string}
//	@Router			/checkout [post]
func (c *checkoutHandlers) Checkout(ctx echo.Context) error {

	req := models.CheckoutRequest{}

	// Bind data from request
	err := ctx.Bind(&req)
	if err != nil {
		return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(err.Error(), http.StatusBadRequest))

	}

	// Validate struct
	err = _validator.Validate().Struct(req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			message := _validator.GetValidatorErrMsg(err.(validator.ValidationErrors))
			return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(message, http.StatusBadRequest))
		}
	}

	id := ctx.Get("id").(int)
	order, err := c.Options.Services.Checkout.Checkout(ctx.Request().Context(), models.Checkout{
		PaymentType: req.PaymentType,
		CustomerID:  id,
	})
	if err != nil {
		return utils.RespondWithError(ctx, customerror.GetStatusCode(err), utils.GetErrorResponse(err.Error(), customerror.GetStatusCode(err)))
	}

	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, order)
}
