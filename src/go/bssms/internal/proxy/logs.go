package proxy

import (
	"bssms/internal/common"
	"log/slog"
)

func getLogger() *slog.Logger {
	return common.GetLoggerFor("proxy")
}
