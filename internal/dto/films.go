package dto

import "task/pkg/database/entities/entitiescustom"

type FilmSearchRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

type FilmsPaginated struct {
	Films    []Film `json:"films"`
	Count    int    `json:"count"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type Film struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type FilmDetail struct {
	ID          int                        `json:"id"`
	Title       string                     `json:"title"`
	Director    string                     `json:"director"`
	ReleaseDate entitiescustom.ReleaseDate `json:"release_date"`
	Synopsis    string                     `json:"synopsis"`
}

type FilmCreateRequest struct {
	Title       string                     `json:"title" validate:"required"`
	Director    string                     `json:"director"`
	ReleaseDate entitiescustom.ReleaseDate `json:"release_date"`
	Synopsis    string                     `json:"synopsis"`
	UserID      int                        `json:"-"`
}

type FilmUpdateRequest struct {
	ID          int                        `param:"id" validate:"required"`
	Title       string                     `json:"title" validate:"required"`
	Director    string                     `json:"director"`
	ReleaseDate entitiescustom.ReleaseDate `json:"release_date"`
	Synopsis    string                     `json:"synopsis"`
	UserID      int                        `json:"-"`
}
