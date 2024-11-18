package handlers

import (
	"matching-service/bin/middlewares"
	driver "matching-service/bin/modules/driver"
	"matching-service/bin/modules/driver/models"

	"matching-service/bin/pkg/utils"

	"github.com/labstack/echo/v4"
)

type driverHttpHandler struct {
	driverUsecaseQuery   driver.UsecaseQuery
	driverUseCaseCommand driver.UsecaseCommand
}

func InitDriverHttpHandler(e *echo.Echo, uq driver.UsecaseQuery, uc driver.UsecaseCommand) {

	handler := &driverHttpHandler{
		driverUsecaseQuery:   uq,
		driverUseCaseCommand: uc,
	}
	route := e.Group("/driver")
	route.POST("/v1/pickup-passanger", handler.PickupPassanger, middlewares.VerifyBearer)
	route.POST("/v1/complete-trip", handler.CompletedTrip, middlewares.VerifyBearer)

}

func (u driverHttpHandler) PickupPassanger(c echo.Context) error {
	var request models.PickupPassanger
	if err := c.Bind(&request); err != nil {
		return utils.ResponseError(err, c)
	}

	if err := request.Validate(); err != nil {
		return utils.ResponseError(err, c)
	}

	userId := utils.ConvertString(c.Get("userId"))
	result := u.driverUseCaseCommand.PickupPassanger(c.Request().Context(), userId, request)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "pickup passanger", 200, c)
}

func (u driverHttpHandler) CompletedTrip(c echo.Context) error {
	req := new(models.Trip)

	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(err, c)
	}

	if err := c.Validate(req); err != nil {
		return utils.ResponseError(err, c)
	}

	userId := utils.ConvertString(c.Get("userId"))
	result := u.driverUseCaseCommand.CompletedTrip(c.Request().Context(), userId, *req)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "completed trip", 200, c)
}
