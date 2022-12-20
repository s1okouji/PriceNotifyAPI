package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/apps", setApp)
	e.Logger.Fatal(e.Start(":1323"))
}

func setApp(c echo.Context) error {

	app := c.FormValue("app_id")
	user := c.FormValue("user_id")

	return c.String(http.StatusOK, "appId:"+app+", userId:"+user)
}
