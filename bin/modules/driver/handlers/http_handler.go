package handlers

import (
	"matching-service/bin/middlewares"
	driver "matching-service/bin/modules/driver"
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
	route.GET("/v1/detail-trip", handler.DetailTrip, middlewares.VerifyBearer)

}

func (u driverHttpHandler) DetailTrip(c echo.Context) error {

	return utils.Response("", "update beacon", 200, c)
}
