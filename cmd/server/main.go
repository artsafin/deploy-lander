package main

import (
	"deploy-lander/internal/application"
	"deploy-lander/internal/config"
	"deploy-lander/internal/data"
	"deploy-lander/internal/httpui"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	fmt.Println("Deploy Lander 1.0")

	projects := config.LoadFromEnv()

	logger, logerr := zap.NewDevelopment(zap.AddCaller(), zap.AddStacktrace(zapcore.DebugLevel))
	if logerr != nil {
		panic(logerr)
	}
	defer logger.Sync()

	builds := data.NewBuildsData()

	app := application.NewApp(projects, logger, builds)
	ui := httpui.NewHttpUI(app, logger)
	ui.ListenAndServe()
}