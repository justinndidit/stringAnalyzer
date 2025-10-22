package application

import (
	"github.com/justinndidit/stringAnalyzer/internal/config"
	"github.com/justinndidit/stringAnalyzer/internal/database"
	"github.com/justinndidit/stringAnalyzer/internal/handler"
	"github.com/justinndidit/stringAnalyzer/internal/repository"
	"github.com/rs/zerolog"
)

type Application struct {
	Config  *config.Config
	Logger  *zerolog.Logger
	DB      *database.Database
	Handler *handler.StringAnalyzerHandler
	repo    *repository.StringRepository
}

func NewApp(cfg *config.Config, logger *zerolog.Logger, db *database.Database) *Application {
	repo := repository.NewStringRepository(logger, db)
	handler := handler.NewStringAnalyzerHandler(logger, db, repo)
	return &Application{
		Config:  cfg,
		Logger:  logger,
		DB:      db,
		repo:    repo,
		Handler: handler,
	}
}
