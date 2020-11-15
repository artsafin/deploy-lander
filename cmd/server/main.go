package main

import (
	"deploy-lander/internal/application"
	"deploy-lander/internal/config"
	"deploy-lander/internal/data"
	"deploy-lander/internal/httpui"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	logger, logerr := zap.NewDevelopment(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if logerr != nil {
		panic(logerr)
	}
	defer logger.Sync()

	fmt.Println("Deploy Lander 1.0")

	projects := config.LoadFromEnv()

	if len(projects) == 0 {
		logger.Fatal("projects configuration is empty")
		os.Exit(1)
	}

	builds := data.NewBuildsData()

	app := application.NewApp(projects, logger, builds)
	ui := httpui.NewHttpUI(app, logger)
	ui.ListenAndServe()
}
