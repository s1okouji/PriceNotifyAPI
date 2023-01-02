package controller

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func PostInteractions(c echo.Context) error {
	var publicKey, _ = hex.DecodeString(os.Getenv("PUBLIC_KEY"))
	var sig, _ = hex.DecodeString(c.Request().Header.Get("X-Signature-Ed25519"))
	var timestamp = c.Request().Header.Get("X-Signature-Timestamp")
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid request body")
	}
	var isVerified = ed25519.Verify(publicKey, []byte(timestamp+string(body)), sig)

	if !isVerified {
		return c.String(http.StatusUnauthorized, "invalid request signature")
	}
	fmt.Println(string(body))
	return c.String(http.StatusOK, "OK")
}
