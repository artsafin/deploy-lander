package httpui

import (
	"go.uber.org/zap"
	"net/http"
)

type errorWriter struct {
	resp   http.ResponseWriter
	req    *http.Request
	logger *zap.Logger
}

func (errw *errorWriter) writeWithError(status int, summary string, err error) {
	errw.logger.Warn(
		summary,
		zap.String("method", errw.req.Method),
		zap.String("path", errw.req.URL.Path),
		zap.Error(err),
	)

	errw.resp.WriteHeader(status)
	errw.resp.Write([]byte(summary))
}

func (errw *errorWriter) write(status int, summary string) {
	errw.logger.Warn(
		summary,
		zap.String("method", errw.req.Method),
		zap.String("path", errw.req.URL.Path),
	)

	errw.resp.WriteHeader(status)
	errw.resp.Write([]byte(summary))
}
