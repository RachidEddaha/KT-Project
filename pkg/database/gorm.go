package database

import (
	"KTOnlinePlatform/pkg/utils"
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewGormDatabase(sqlDB *sql.DB) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		FullSaveAssociations: true,
		Logger:               gormLogger.Default.LogMode(gormLogger.Info),
		NowFunc:              utils.TimeNowInUTC,
		TranslateError:       true,
	}
	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB}),
		gormConfig,
	)
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
