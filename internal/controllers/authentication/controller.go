package authentication

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"task/internal/dto"
	"task/pkg/logger"
)

type service interface {
	Login(ctx context.Context, request dto.Login) (dto.JWTTokens, error)
	CreateUser(ctx context.Context, request dto.CreateUserRequest) error
}

type Controller struct {
	service service
}

func NewController(service service) *Controller {
	if service == nil {
		panic(service)
	}
	return &Controller{
		service: service,
	}
}

func (c *Controller) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/api/v1")

	g.PUT("/login", c.logIn)
	g.POST("/register", c.createUser)
	// should exist in the future, to refresh the token when the access token is expired
	// g.POST("/refresh-token", c.refresh) // TODO: implement refresh token
}

func (c *Controller) logIn(context echo.Context) error {
	request := dto.Login{}
	err := context.Bind(&request)
	if err != nil {
		return err
	}
	err = context.Validate(request)
	if err != nil {
		logger.Error().Err(err).Msg("validation failed")
		return err
	}

	tokens, err := c.service.Login(context.Request().Context(), request)
	if err != nil {
		logger.Error().Err(err).Msg("login failed")
		return err
	}
	return context.JSON(http.StatusOK, tokens)
}

func (c *Controller) createUser(context echo.Context) error {
	request := dto.CreateUserRequest{}
	err := context.Bind(&request)
	if err != nil {
		return err
	}
	err = context.Validate(request)
	if err != nil {
		logger.Error().Err(err).Msg("validation failed")
		return err
	}

	err = c.service.CreateUser(context.Request().Context(), request)
	if err != nil {
		logger.Error().Err(err).Msg("create user failed")
		return err
	}
	return context.NoContent(http.StatusOK)
}
