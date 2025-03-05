package main

import (
	authcontroller "task/internal/controllers/authentication"
	"task/internal/repositories/authentication"
	authservice "task/internal/services/authentication"
	"task/pkg/configuration"
	"task/pkg/database"
	"task/pkg/logger"
	"task/pkg/middlewares"
	"task/pkg/webutils"
)

func main() {
	config, err := configuration.LoadConfiguration()
	if err != nil {
		panic(err)
	}
	logger.Initialize(config.ConfigLogger)
	db := database.NewDatabase(config.ConfigDatabase)
	e := webutils.NewEcho(config.ConfigEcho)

	middleware := middlewares.NewMiddleware(config.JWTSecret)
	authRep := authentication.NewRepository(db)
	authService := authservice.NewService(authRep, middleware)
	authcontroller.NewController(authService).RegisterRoutes(e)

	webutils.StartEcho(e, config.AddressEcho)
}
