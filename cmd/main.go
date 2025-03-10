package main

import (
	authcontroller "KTOnlinePlatform/internal/controllers/authentication"
	filmscontroller "KTOnlinePlatform/internal/controllers/films"
	"KTOnlinePlatform/internal/repositories/authentication"
	"KTOnlinePlatform/internal/repositories/films"
	authservice "KTOnlinePlatform/internal/services/authentication"
	filmsservice "KTOnlinePlatform/internal/services/films"
	"KTOnlinePlatform/pkg/configuration"
	"KTOnlinePlatform/pkg/database"
	"KTOnlinePlatform/pkg/logger"
	"KTOnlinePlatform/pkg/middlewares"
	"KTOnlinePlatform/pkg/webutils"
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
