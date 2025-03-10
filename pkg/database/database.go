package database

import (
	"KTOnlinePlatform/pkg/configuration"
	"KTOnlinePlatform/pkg/logger"
	"gorm.io/gorm"
)

func NewDatabase(config configuration.ConfigDatabase) *gorm.DB {
	sqlDB, err := NewSqlDB(config)
	if err != nil {
		logger.Fatal().Msgf("Error while connecting to database: %v", err)
	}
	gormDB, err := NewGormDatabase(sqlDB)
	if err != nil {
		logger.Fatal().Msgf("Error while connecting to database: %v", err)
	}
	return gormDB
}
