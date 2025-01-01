package installer

import (
	"bssms/internal/common"
	"log/slog"
)

func getLogger() *slog.Logger {
	return common.GetLoggerFor("installer")
}
