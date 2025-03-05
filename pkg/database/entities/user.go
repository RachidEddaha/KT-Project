package entities

import (
	"time"
)

type User struct {
	ID        int        `db:"id"  json:"id"`
	Username  string     `db:"username" json:"username"`
	Password  string     `db:"password" json:"password"`
	CreatedAt *time.Time `db:"created_at" gorm:"column:created_at;type:TIMESTAMPTZ;" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" gorm:"column:updated_at;type:TIMESTAMPTZ;" json:"updatedAt"`
}
