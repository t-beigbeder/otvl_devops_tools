package common

import "log/slog"

func GetLogger() *slog.Logger {
	return slog.Default()
}

func GetLoggerFor(app string) *slog.Logger {
	return GetLogger().With("app", app)
}
