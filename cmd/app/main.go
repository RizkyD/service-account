package main

import (
	"account-service/internal/config"
	"account-service/internal/handler"
	"account-service/internal/middleware"
	"account-service/internal/repository"
	"account-service/internal/route"
	"account-service/internal/service"
	"account-service/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	//setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zlog.Info().Msg("Starting application setup...")

	//setup config
	configPath := "./config.json"
	cfg := config.LoadConfig(configPath)

	//setup database
	dbPool := database.ConnectDB(cfg.DatabaseURL)
	defer dbPool.Close()

	//setup depedency injection
	nasabahRepository := repository.NewNasabahRepository(dbPool)
	nasabahService := service.NewNasabahService(nasabahRepository, dbPool)
	nasabahHandler := handler.NewNasabahHandler(nasabahService)

	//setup fiber
	app := fiber.New()

	//use middleware
	app.Use(middleware.ZerologRequestLogger)

	//setup router
	route.SetupRoutes(app, nasabahHandler)

	zlog.Info().Str("port", cfg.ServerPort).Msgf("Starting server '%s'", cfg.AppName)
	err := app.Listen(cfg.ServerPort)
	if err != nil {
		zlog.Fatal().Err(err).Msg("FATAL: Could not start server")
	}
}
