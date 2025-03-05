package database

import (
	"database/sql"
	"task/pkg/configuration"
	"task/pkg/logger"
	"time"
)

const (
	retryDelay          = 3 * time.Second
	maxAttemptToConnect = 5
)

func NewSqlDB(config configuration.ConfigDatabase) (*sql.DB, error) {
	logger.Debug().Msgf("Trying to connect to database with config %s...", config)

	db, err := sql.Open("pgx", config.GetDSN())
	if err != nil {
		logger.Error().Msgf("Error while connecting to database: %v", err)
		return nil, err
	}

	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)

	for i := 0; i < maxAttemptToConnect; i++ {
		err = db.Ping()
		if err != nil {
			logger.Warn().Msgf("%v - retry in %v, attempt: %d out of %d", err, retryDelay, i, maxAttemptToConnect)
			time.Sleep(retryDelay)
			continue
		}
		break
	}

	if err != nil {
		logger.Error().Msgf("Cannot connect to database: %v", err)
		return nil, err
	}

	logger.Info().Msgf("Connected to database with config: %s", config)
	return db, nil
}
