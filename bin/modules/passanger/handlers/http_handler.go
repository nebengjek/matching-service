package handlers

import (
	"matching-service/bin/middlewares"
	driver "matching-service/bin/modules/passanger"
	httpError "matching-service/bin/pkg/http-error"
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
	route.GET("/v1/trip-request", handler.DetailTrip, middlewares.VerifyBearer)

}

func (u driverHttpHandler) DetailTrip(c echo.Context) error {
	passangerId := c.QueryParam("psgId")
	if passangerId == "" {
		errObj := httpError.BadRequest("need params")
		return utils.ResponseError(errObj, c)
	}

	userId := utils.ConvertString(c.Get("userId"))
	result := u.driverUsecaseQuery.DetailTrip(c.Request().Context(), userId, passangerId)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "update beacon", 200, c)
}
