package main

import (
	"github.com/labstack/echo/v4"
	"github.com/s1okouji/price_notify_api/controller"
	"github.com/s1okouji/price_notify_api/service"
)

func main() {
	service.SetUp()
	e := echo.New()
	e.GET("/apps", controller.GetApps)
	e.POST("/apps", controller.SetApp)
	e.DELETE("/apps", controller.DeleteApp)
	e.Logger.Fatal(e.Start(":8000"))
}
