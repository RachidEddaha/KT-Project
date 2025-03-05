package main

import (
	authcontroller "task/internal/controllers/authentication"
	filmscontroller "task/internal/controllers/films"
	"task/internal/repositories/authentication"
	"task/internal/repositories/films"
	authservice "task/internal/services/authentication"
	filmsservice "task/internal/services/films"
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

	filmRepo := films.NewRepository(db)
	filmService := filmsservice.NewService(filmRepo)
	filmscontroller.NewController(filmService, middleware).RegisterRoutes(e)

	webutils.StartEcho(e, config.AddressEcho)
}
