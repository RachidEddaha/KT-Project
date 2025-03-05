package films

import (
	"context"
	"github.com/samber/lo"
	"task/internal/dto"
	"task/internal/models"
	"task/internal/models/consts"
	"task/internal/models/kterrors"
	"task/pkg/customerror"
	"task/pkg/database/entities"
)

type Repository interface {
	GetFilmsPaginated(ctx context.Context, pageSize int, offset int) ([]models.FilmPaginated, error)
	GetFilm(ctx context.Context, ID int) (entities.Film, error)
	DeleteFilm(ctx context.Context, id int) error
	CreateFilm(ctx context.Context, film entities.Film) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetFilmPaginated(ctx context.Context, request dto.FilmSearchRequest) (dto.FilmsPaginated, error) {
	offset := calculateOffset(request.Page, request.PageSize)
	result, err := s.repo.GetFilmsPaginated(ctx, request.PageSize, offset)
	if err != nil {
		return dto.FilmsPaginated{}, err
	}
	if len(result) == 0 {
		return dto.FilmsPaginated{
			Page:     request.Page,
			PageSize: request.PageSize,
		}, nil
	}
	films := lo.Map(result, func(item models.FilmPaginated, index int) dto.Film {
		return dto.Film{
			ID:    item.ID,
			Title: item.Title,
		}
	})
	return dto.FilmsPaginated{
		Films:    films,
		Count:    result[0].Qty,
		Page:     request.Page,
		PageSize: request.PageSize,
	}, nil
}

func calculateOffset(page, size int) int {
	offset := (page - 1) * size
	if offset < 0 {
		offset = consts.BasicPaginationDefaultOffset
	}
	return offset
}

func (s *Service) GetFilmDetail(ctx context.Context, ID int) (dto.FilmDetail, error) {
	film, err := s.repo.GetFilm(ctx, ID)
	if err != nil {
		return dto.FilmDetail{}, err
	}
	return dto.FilmDetail{
		ID:          film.ID,
		Title:       film.Title,
		Director:    film.Director,
		ReleaseDate: film.ReleaseDate,
		Synopsis:    film.Synopsis,
	}, nil
}

func (s *Service) DeleteFilm(ctx context.Context, filmID int, userID int) error {
	film, err := s.repo.GetFilm(ctx, filmID)
	if err != nil {
		return err
	}
	if film.UserID != userID {
		return customerror.NewCustomError(kterrors.UserCannotDeleteFilmError)
	}
	return s.repo.DeleteFilm(ctx, filmID)
}

func (s *Service) CreateFilm(ctx context.Context, request dto.FilmCreateRequest) error {
	film := entities.Film{
		Title:       request.Title,
		Director:    request.Director,
		ReleaseDate: request.ReleaseDate,
		Synopsis:    request.Synopsis,
		UserID:      request.UserID,
	}
	err := s.repo.CreateFilm(ctx, film)
	if err != nil {
		if customerror.IsUniqueViolation(err) {
			return customerror.NewCustomError(kterrors.FilmTitleAlreadyExistsError)
		}
		return err
	}
	return nil
}
