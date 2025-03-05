package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

func TimeNowInUTC() time.Time {
	return time.Now().In(time.UTC)
}

func GetUserID(e echo.Context) (int, error) {
	token := e.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	s := claims["sub"]
	sub, err := strconv.Atoi(s.(string))
	if err != nil {
		return 0, err
	}
	return sub, nil
}
