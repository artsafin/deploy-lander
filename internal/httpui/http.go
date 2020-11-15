package httpui

import (
	"deploy-lander/internal/application"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type httpUI struct {
	app    *application.App
	logger *zap.Logger
}

func NewHttpUI(app *application.App, logger *zap.Logger) *httpUI {
	return &httpUI{app, logger}
}

func (ui *httpUI) ListenAndServe() {
	http.HandleFunc("/build", func(wr http.ResponseWriter, req *http.Request) {
		errw := errorWriter{resp: wr, req: req, logger: ui.logger}
		if req.Method != http.MethodPost {
			errw.write(http.StatusBadRequest, "Only POST method is supported")
			return
		}

		alias := strings.Trim(req.URL.Path, "/")

		if len(alias) == 0 {
			errw.write(http.StatusBadRequest, "Alias must be provided")
			return
		}

		version, parseerr := parseVersion(req)
		if parseerr != nil {
			errw.writeWithError(http.StatusBadRequest, "Invalid parameters", parseerr)
			return
		}

		err := ui.app.RegisterBuild(alias, version)
		if err != nil {
			errw.writeWithError(http.StatusBadRequest, fmt.Sprint(err), err)
			return
		}

		ui.logger.Info("OK RegisterBuild", zap.String("path", req.URL.Path))
		wr.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/deploy", func(wr http.ResponseWriter, req *http.Request) {
		errw := errorWriter{resp: wr, req: req, logger: ui.logger}
		if req.Method != http.MethodPost {
			errw.write(http.StatusBadRequest, "Only POST method is supported")
			return
		}

		alias := strings.Trim(req.URL.Path, "/")

		if len(alias) == 0 {
			errw.write(http.StatusBadRequest, "Alias must be provided")
			return
		}

		version, parseerr := parseVersion(req)
		if parseerr != nil {
			errw.writeWithError(http.StatusBadRequest, "Invalid parameters", parseerr)
			return
		}

		err := ui.app.DoDeploy(alias, version)
		if err != nil {
			errw.writeWithError(http.StatusBadRequest, fmt.Sprint(err), err)
			return
		}

		ui.logger.Info("OK DoDeploy", zap.String("path", req.URL.Path))
		wr.WriteHeader(http.StatusOK)
	})

	httperr := http.ListenAndServe(":8080", nil)
	ui.logger.Fatal("HTTP", zap.Error(httperr))
}
