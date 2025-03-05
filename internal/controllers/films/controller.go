package authentication

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"task/internal/dto"
	"task/internal/models/consts"
	"task/pkg/logger"
	"task/pkg/middlewares"
	"task/pkg/utils"
	"task/pkg/webutils"
)

type service interface {
	GetFilmPaginated(ctx context.Context, request dto.FilmSearchRequest) (dto.FilmsPaginated, error)
	GetFilmDetail(ctx context.Context, ID int) (dto.FilmDetail, error)
	DeleteFilm(ctx context.Context, filmID int, userID int) error
	CreateFilm(ctx context.Context, request dto.FilmCreateRequest) error
}

type Controller struct {
	service service
	middlewares.AuthMiddleware
}

func NewController(service service, middleware middlewares.AuthMiddleware) *Controller {
	if service == nil {
		panic(service)
	}
	if middleware == nil {
		panic(middleware)
	}
	return &Controller{
		service:        service,
		AuthMiddleware: middleware,
	}
}

func (c *Controller) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/api/v1/films", c.AuthMiddleware.Authenticated())

	g.GET("", c.getFilmPaginated) // TODO: include query params to filter the films
	g.GET("/:id", c.getFilmDetail)
	// g.PUT("/:id", c.updateFilmDetail)
	g.DELETE("/:id", c.deleteFilm)
	g.POST("", c.createFilm)
}

func (c *Controller) getFilmPaginated(context echo.Context) error {
	request := dto.FilmSearchRequest{}
	err := context.Bind(&request)
	if err != nil {
		return err
	}
	if request.Page == 0 {
		request.Page = consts.BasicPaginationDefaultPageNumber
	}

	if request.PageSize == 0 {
		request.PageSize = consts.PaginationDefaultPageSize
	}

	result, err := c.service.GetFilmPaginated(context.Request().Context(), request)
	if err != nil {
		logger.Error().Err(err).Msg("get film paginated failed")
		return err
	}
	return context.JSON(http.StatusOK, result)
}

func (c *Controller) getFilmDetail(context echo.Context) error {
	filmID, err := webutils.CheckParamToInt(context, "id")
	if err != nil {
		return err
	}

	result, err := c.service.GetFilmDetail(context.Request().Context(), filmID)
	if err != nil {
		logger.Error().Err(err).Msg("get film detail failed")
		return err
	}
	return context.JSON(http.StatusOK, result)
}

func (c *Controller) deleteFilm(context echo.Context) error {
	filmID, err := webutils.CheckParamToInt(context, "id")
	if err != nil {
		return err
	}

	userID, err := utils.GetUserID(context)
	if err != nil {
		return err
	}

	err = c.service.DeleteFilm(context.Request().Context(), filmID, userID)
	if err != nil {
		logger.Error().Err(err).Msg("delete film failed")
		return err
	}
	return context.NoContent(http.StatusOK)
}

func (c *Controller) createFilm(context echo.Context) error {
	request := dto.FilmCreateRequest{}
	err := context.Bind(&request)
	if err != nil {
		return err
	}
	err = context.Validate(request)
	if err != nil {
		return err
	}
	request.UserID, err = utils.GetUserID(context)
	if err != nil {
		return err
	}

	err = c.service.CreateFilm(context.Request().Context(), request)
	if err != nil {
		logger.Error().Err(err).Msg("create film failed")
		return err
	}
	return context.NoContent(http.StatusOK)
}
