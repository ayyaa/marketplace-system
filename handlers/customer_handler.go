package handlers

import (
	"marketplace-system/lang"
	"marketplace-system/models"
	"net/http"
	"strings"

	customerror "marketplace-system/lib/customerrors"
	utils "marketplace-system/lib/helper"
	_validator "marketplace-system/lib/validator"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type customerHandlers handlers

type CustomerInterface interface {
	InsertCustomer(ctx echo.Context) error
	Login(ctx echo.Context) error
}

// Customer godoc
//
//	@Summary		Create Customer / Register
//	@Description	Create customer / register
//	@ID				customer-create
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CustomerRequest	true	"Create customer"
//	@Success		201		{object}	models.Response{data=[]models.Customer}
//	@Failure		400		{object}	models.BasicResponse{message=[]string}
//	@Failure		500		{object}	models.BasicResponse{message=[]string}
//	@Router			/register [post]
func (h *customerHandlers) InsertCustomer(ctx echo.Context) error {

	req := models.CustomerRequest{}

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
	// call method insert to database
	_, err = h.Options.Services.Customer.InsertCustomer(ctx.Request().Context(), req)
	if err != nil {
		// check if error contains "customer_email_key", phone number already exist
		if strings.Contains(err.Error(), "customer_phone_key") {
			return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(lang.ErrorConflictPhoneNumber, http.StatusBadRequest))
		}

		if strings.Contains(err.Error(), "customer_email_key") {
			return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(lang.ErrorConflictEmail, http.StatusBadRequest))
		}

		return utils.RespondWithError(ctx, http.StatusInternalServerError, utils.GetErrorResponse(err.Error(), http.StatusInternalServerError))
	}

	// return success
	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, nil)
}

// Customer godoc
//
//	@Summary		Login
//	@Description	Login
//	@ID				customer-login
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CustomerRequest	true	"Create customer"
//	@Success		201		{object}	models.Response{data=[]models.Customer}
//	@Failure		400		{object}	models.BasicResponse{message=[]string}
//	@Failure		500		{object}	models.BasicResponse{message=[]string}
//	@Router			/register [post]
func (c *customerHandlers) Login(ctx echo.Context) error {

	// intial
	var req models.LoginRequest

	// Bind data from request
	err := ctx.Bind(&req)
	if err != nil {
		return utils.RespondWithError(ctx, http.StatusBadRequest, utils.GetErrorResponse(err.Error(), http.StatusBadRequest))
	}

	// validate if any user by phone
	token, err := c.Options.Services.Customer.Login(ctx.Request().Context(), req)
	if err != nil {
		return utils.RespondWithError(ctx, customerror.GetStatusCode(err), utils.GetErrorResponse(err.Error(), customerror.GetStatusCode(err)))
	}

	// return error
	logrus.Info(lang.SuccessMsg)
	return utils.ResponseSuccess(ctx, models.LoginResponse{
		Token: token,
	})

}
