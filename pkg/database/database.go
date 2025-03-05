package database

import (
	"gorm.io/gorm"
	"task/pkg/configuration"
	"task/pkg/logger"
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
