package handlers

import (
	"marketplace-system/lang"
	customerror "marketplace-system/lib/customerrors"
	utils "marketplace-system/lib/helper"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type productHandlers handlers

type ProductInterface interface {
	FindProductsByCategory(ctx echo.Context) error
}

// Customer godoc
//
//	@Summary		Get Product by Categories
//	@Description	Get Product by Categories
//	@ID				get-product-by-categories
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"Param slug"
//	@Success		201		{object}	models.Response{data=[]models.Customer}
//	@Failure		400		{object}	models.BasicResponse{message=[]string}
//	@Failure		500		{object}	models.BasicResponse{messsage=[]string}
//	@Router			/categories/{slug}/products [get]
func (h *productHandlers) FindProductsByCategory(ctx echo.Context) error {
	categorySlug := ctx.Param("slug")

	products, err := h.Options.Services.Product.FindProductsByCategory(ctx.Request().Context(), categorySlug)
	if err != nil {
		return utils.RespondWithError(ctx, customerror.GetStatusCode(err), utils.GetErrorResponse(err.Error(), customerror.GetStatusCode(err)))
	}

	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, products)
}
