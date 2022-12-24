package controller

import (
	"encoding/json"
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
	return c.String(http.StatusCreated, "OK")
}

func DeleteApp(c echo.Context) error {
	var dto dto.DeleteAppDTO
	err := c.Bind(&dto)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	service.DeleteGame(dto)
	return c.String(http.StatusOK, "OK")
}

func GetApps(c echo.Context) error {
	apps := service.GetGames()
	bytes, _ := json.Marshal(&apps)
	return c.String(http.StatusOK, string(bytes))
}
