package application

import (
	"deploy-lander/internal/data"
	"deploy-lander/internal/structs"
	"go.uber.org/zap"
)

type App struct {
	projects structs.ProjectMap
	logger *zap.Logger
	builds *data.BuildsData
}

func NewApp(projects structs.ProjectMap, logger *zap.Logger, builds *data.BuildsData) *App {
	return &App{projects, logger, builds}
}
