package entities

import (
	"task/pkg/database/entities/entitiescustom"
	"time"
)

type Film struct {
	ID          int                        `db:"id"  json:"id"`
	Title       string                     `db:"title" json:"title"`
	Director    string                     `db:"director" json:"director"`
	ReleaseDate entitiescustom.ReleaseDate `db:"release_date" json:"release_date"`
	Synopsis    string                     `db:"synopsis" json:"synopsis"`
	UserID      int                        `db:"user_id" json:"user_id"`
	CreatedAt   *time.Time                 `db:"created_at" gorm:"column:created_at;type:TIMESTAMPTZ;" json:"createdAt"`
	UpdatedAt   *time.Time                 `db:"updated_at" gorm:"column:updated_at;type:TIMESTAMPTZ;" json:"updatedAt"`
}
