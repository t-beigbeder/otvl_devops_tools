package barsvc

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/t-beigbeder/otvl_devops_tools/src/go/bar_service/svcctl"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
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

type barService struct {
	address, backup, restore           string
	currentStatus, lastOperationStatus string
	lastOperationDate                  time.Time
	sync                               sync.Mutex
	e                                  *echo.Echo
}

func (bs *barService) configFromEnv() {
	bs.address = ":3000"
	bs.backup = "/bin/sh -c /etc/bar/backup.sh"
	bs.restore = "/bin/sh -c /etc/bar/restore.sh"
	if os.Getenv("BAR_ADDRESS") != "" {
		bs.address = os.Getenv("BAR_ADDRESS")
	}
	if os.Getenv("BAR_BACKUP") != "" {
		bs.backup = os.Getenv("BAR_BACKUP")
	}
	if os.Getenv("BAR_RESTORE") != "" {
		bs.restore = os.Getenv("BAR_RESTORE")
	}
}

func (bs *barService) configure(address, backup, restore string) {
	bs.address, bs.backup, bs.restore = address, backup, restore
}

func (bs *barService) logOutErr(out, err string) {
	if out != "" {
		bs.Logger().Info(out)
	}
	if err != "" {
		bs.Logger().Info(err)
	}
}

func (bs *barService) bor(c echo.Context, isRestore bool) error {
	sCmd := bs.backup
	if isRestore {
		sCmd = bs.restore
	}
	args := strings.Split(sCmd, " ")
	cmd := exec.Command(args[0], args[1:]...)
	var out, ser strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &ser
	if err := cmd.Run(); err != nil {
		bs.logOutErr(out.String(), ser.String())
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "OK")
}

func (bs *barService) Name() string {
	return "BackupAndRestoreService"
}

func (bs *barService) Start() error {
	bs.Logger().SetLevel(log.INFO)
	bs.configFromEnv()
	e := bs.e
	e.POST("/backup", func(c echo.Context) error {
		return bs.bor(c, false)
	})
	e.POST("/restore", func(c echo.Context) error {
		return bs.bor(c, true)
	})
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	if err := e.Start(bs.address); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
	e.Logger.Info("Start done")
	return nil
}

func (bs *barService) Stop(ctx context.Context) error {
	return bs.e.Shutdown(ctx)
}

func (bs *barService) Logger() echo.Logger {
	return bs.e.Logger
}

func newSvc() svcctl.ControllableService {
	bs := &barService{e: echo.New()}
	bs.configFromEnv()
	return bs
}

func BarSvc() svcctl.ControllableService {
	return newSvc()
}
