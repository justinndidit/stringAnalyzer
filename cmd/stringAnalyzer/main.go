package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/justinndidit/stringAnalyzer/internal/application"
	"github.com/justinndidit/stringAnalyzer/internal/config"
	"github.com/justinndidit/stringAnalyzer/internal/database"
	"github.com/justinndidit/stringAnalyzer/internal/logger"
	"github.com/justinndidit/stringAnalyzer/internal/routes"
	"github.com/justinndidit/stringAnalyzer/internal/server"
)

const DefaultContextTimeout = 30

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logger := logger.NewLogger()

	if err = database.Migrate(context.Background(), &logger, cfg); err != nil {
		logger.Fatal().Err(err).Msg("failed to migrate database")
	}

	db, err := database.New(cfg, &logger)

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize database")
	}

	app := application.NewApp(cfg, &logger, db)

	r := routes.SetupAuthRoutes(app)

	srv, err := server.New(app, cfg)

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize server")
	}

	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Start server
	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server forced to shutdown")
	}
	stop()
	cancel()

	logger.Info().Msg("server exited properly")

}
