package handlers

import (
	"marketplace-system/lang"
	customerror "marketplace-system/lib/customerrors"
	utils "marketplace-system/lib/helper"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type orderHandlers handlers

type OrderInterface interface {
	GetOrders(ctx echo.Context) error
	GetOrdersById(ctx echo.Context) error
}

// Order godoc
//
//	@Summary		Get orders
//	@Description	Get orders
//	@ID				get-orders
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Success		201		{object}	models.ResponseSuccess{data=models.Order}
//	@Failure		400		{object}	models.ApplicationError{message=[]string}
//	@Failure		500		{object}	models.ApplicationError{messsage=[]string}
//	@Router			/orders [post]
func (c *orderHandlers) GetOrders(ctx echo.Context) error {

	id := ctx.Get("id").(int)
	order, err := c.Options.Services.Checkout.GetOrders(ctx.Request().Context(), id)
	if err != nil {
		return utils.RespondWithError(ctx, customerror.GetStatusCode(err), utils.GetErrorResponse(err.Error(), customerror.GetStatusCode(err)))
	}

	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, order)
}

// Order godoc
//
//	@Summary		Get order by invoice Id
//	@Description	Get orders  by invoice Id
//	@ID				get-orders-by-invoice-id
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Success		201		{object}	models.ResponseSuccess{data=models.Order}
//	@Failure		400		{object}	models.ApplicationError{message=[]string}
//	@Failure		500		{object}	models.ApplicationError{messsage=[]string}
//	@Router			/order/:invoice [post]
func (c *orderHandlers) GetOrdersById(ctx echo.Context) error {
	invoice := ctx.Param("invoice")
	order, err := c.Options.Services.Checkout.GetOrderByInvoice(ctx.Request().Context(), invoice)
	if err != nil {
		return utils.RespondWithError(ctx, customerror.GetStatusCode(err), utils.GetErrorResponse(err.Error(), customerror.GetStatusCode(err)))
	}

	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, order)
}
