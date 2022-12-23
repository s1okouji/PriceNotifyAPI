package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/s1okouji/price_notify_api/dto"
	"github.com/s1okouji/price_notify_api/service"
)

func SetApp(c echo.Context) error {
	var dto dto.CreateAppDTO
	err := c.Bind(&dto)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	service.AddGame(dto)
	return c.String(http.StatusAccepted, "OK")
}
