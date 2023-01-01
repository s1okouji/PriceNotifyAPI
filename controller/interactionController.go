package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PostInteractions(c echo.Context) error {
	println("Accept Request!")
	return c.String(http.StatusOK, "Interactioned! Hello!\n")
}
