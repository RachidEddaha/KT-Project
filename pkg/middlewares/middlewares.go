package middlewares

import (
	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"strconv"
	"task/internal/dto"
	"task/pkg/logger"
	"task/pkg/utils"
	"time"
)

// global interface since it will be used in multiple places
type AuthMiddleware interface {
	Authenticated() echo.MiddlewareFunc
}

type Middleware struct {
	jwtSecret string
}

func NewMiddleware(jwtSecret string) *Middleware {
	if jwtSecret == "" {
		panic(jwtSecret)
	}

	return &Middleware{
		jwtSecret: jwtSecret,
	}
}

func (m *Middleware) Authenticated() echo.MiddlewareFunc {
	return m.authenticatedWithTokenLookup("")
}

func (m *Middleware) authenticatedWithTokenLookup(tokenLookup string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwtMiddleware := m.configureJWT(tokenLookup)
			setNoCacheHeaders(c)
			if err := jwtMiddleware(next)(c); err != nil {
				return err
			}
			return nil
		}
	}
}

func (m *Middleware) configureJWT(tokenLookup string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(m.jwtSecret),
		TokenLookup: tokenLookup,
	})
}

func setNoCacheHeaders(c echo.Context) {
	c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")
}

func (m *Middleware) GenerateAuthTokens(userID int, username string) (dto.JWTTokens, error) {
	t := createClaims(userID, username, 15)
	accessToken, err := t.SignedString([]byte(m.jwtSecret))
	if err != nil {
		return dto.JWTTokens{}, err
	}

	rt := createClaims(userID, username, 1440)
	refreshToken, err := rt.SignedString([]byte(m.jwtSecret))
	if err != nil {
		return dto.JWTTokens{}, err
	}

	logger.Debug().Msgf("Created jwt tokens for %q", username)
	return dto.JWTTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func createClaims(userID int, username string, validationTimeInMinutes int) *jwt.Token {
	t := jwt.New(jwt.SigningMethodHS256)
	tClaims := t.Claims.(jwt.MapClaims)
	tClaims["sub"] = strconv.Itoa(userID)
	tClaims["username"] = username
	tClaims["exp"] = utils.TimeNowInUTC().Add(time.Duration(validationTimeInMinutes) * time.Minute).Unix()

	return t
}
