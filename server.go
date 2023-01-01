package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/s1okouji/price_notify_api/controller"
	"github.com/s1okouji/price_notify_api/service"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	service.SetUp()
	e := echo.New()
	e.AutoTLSManager.Cache = autocert.DirCache("/var/api/.cache")
	e.Pre(middleware.HTTPSRedirect())
	e.GET("/apps", controller.GetApps)
	e.POST("/apps", controller.SetApp)
	e.DELETE("/apps", controller.DeleteApp)
	e.POST("/interactions", controller.PostInteractions)
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}
