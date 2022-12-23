package main

import (
	"github.com/labstack/echo/v4"
	"github.com/s1okouji/price_notify_api/controller"
	"github.com/s1okouji/price_notify_api/service"
)

func main() {
	service.SetUp()
	e := echo.New()
	e.POST("/apps", controller.SetApp)
	e.Logger.Fatal(e.Start(":1323"))
}
