package films

import (
	"KTOnlinePlatform/internal/models"
	"KTOnlinePlatform/pkg/database/entities"
	"context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

const (
	getFilmsPaginated = `
SELECT
		f.id,
		f.title,
			COUNT(*) OVER() AS qty
		FROM films f
		ORDER BY f.title
		LIMIT ? OFFSET ?
`
)

func (r *Repository) GetFilmsPaginated(ctx context.Context, pageSize int, offset int) (result []models.FilmPaginated, err error) {
	err = r.db.WithContext(ctx).Raw(getFilmsPaginated, pageSize, offset).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetFilm(ctx context.Context, ID int) (film entities.Film, err error) {
	err = r.db.WithContext(ctx).First(&film, ID).Error
	if err != nil {
		return film, err
	}
	return film, nil
}

func (r *Repository) DeleteFilm(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&entities.Film{}, id).Error
}

func (r *Repository) CreateFilm(ctx context.Context, film entities.Film) error {
	return r.db.WithContext(ctx).Create(&film).Error
}

func (r *Repository) UpdateFilm(ctx context.Context, film entities.Film) error {
	return r.db.WithContext(ctx).Model(&film).Updates(entities.Film{
		Title:       film.Title,
		Director:    film.Director,
		ReleaseDate: film.ReleaseDate,
		Synopsis:    film.Synopsis,
	}).Error
}
