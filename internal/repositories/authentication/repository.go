package authentication

import (
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

func (r *Repository) FindUser(ctx context.Context, username string) (user entities.User, err error) {
	err = r.db.WithContext(ctx).First(&user, "username = ?", username).Error
	if err != nil {
		return user, err
	}
	if user == (entities.User{}) {
		return user, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, username, password string) error {
	user := entities.User{
		Username: username,
		Password: password,
	}
	err := r.db.WithContext(ctx).Model(&entities.User{}).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
