package barsvc

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

func configFromEnv() (address, backup, restore string) {
	address = ":3000"
	backup = "/bin/sh -c /etc/bar/backup.sh"
	restore = "/bin/sh -c /etc/bar/restore.sh"
	if os.Getenv("BAR_ADDRESS") != "" {
		address = os.Getenv("BAR_ADDRESS")
	}
	if os.Getenv("BAR_BACKUP") != "" {
		address = os.Getenv("BAR_BACKUP")
	}
	if os.Getenv("BAR_RESTORE") != "" {
		address = os.Getenv("BAR_RESTORE")
	}
	return
}

func Service() {
	address, backup, restore := configFromEnv()
	logOutErr := func(c echo.Context, out, err string) {
		if out != "" {
			c.Logger().Printf(out)
		}
		if err != "" {
			c.Logger().Printf(err)
		}
	}
	bar := func(c echo.Context, isRestore bool) error {
		sCmd := backup
		if isRestore {
			sCmd = restore
		}
		args := strings.Split(sCmd, " ")
		cmd := exec.Command(args[0], args[1:]...)
		var out, ser strings.Builder
		cmd.Stdout = &out
		cmd.Stderr = &ser
		if err := cmd.Run(); err != nil {
			logOutErr(c, out.String(), ser.String())
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, "OK")
	}
	// Setup
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.POST("/backup", func(c echo.Context) error {
		return bar(c, false)
	})
	e.POST("/restore", func(c echo.Context) error {
		return bar(c, true)
	})
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
