package barsvc

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/t-beigbeder/otvl_devops_tools/src/go/bar_service/svcctl"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type barService struct {
	address                            string
	backup, restore                    []string
	currentStatus, lastOperationStatus string
	started                            time.Time
	ended                              time.Time
	sync                               sync.Mutex
	e                                  *echo.Echo
	logger                             *log.Logger
}

type mStatus struct {
	Current       string `json:"current,omitempty"`
	LastOperation string `json:"lastOperation,omitempty"`
	Started       string `json:"started,omitempty"`
	Ended         string `json:"ended,omitempty"`
}

type mErr struct {
	msg string `json:"msg,omitempty"`
}

func (bs *barService) configFromEnv() {
	bs.address = ":3000"
	sBackup := "/bin/sh -c /etc/bar/backup.sh"
	sRestore := "/bin/sh -c /etc/bar/restore.sh"
	if os.Getenv("BAR_ADDRESS") != "" {
		bs.address = os.Getenv("BAR_ADDRESS")
	}
	if os.Getenv("BAR_BACKUP") != "" {
		sBackup = os.Getenv("BAR_BACKUP")
	}
	if os.Getenv("BAR_RESTORE") != "" {
		sRestore = os.Getenv("BAR_RESTORE")
	}
	bs.backup = strings.Split(sBackup, " ")
	bs.restore = strings.Split(sRestore, " ")
}

func (bs *barService) configure(address string, backup, restore []string) {
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
	args := bs.backup
	if isRestore {
		args = bs.restore
	}
	bs.sync.Lock()
	locked := false
	if bs.currentStatus == "" {
		bs.currentStatus = fmt.Sprintf("running %s", strings.Join(args, " "))
		bs.started = time.Now()

	} else {
		locked = true
	}
	bs.sync.Unlock()
	if locked {
		return c.JSON(http.StatusOK, &mErr{
			msg: fmt.Sprintf("operation in progress: %s", bs.currentStatus),
		})
	}
	cmd := exec.Command(args[0], args[1:]...)
	var out, ser strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &ser
	if err := cmd.Run(); err != nil {
		bs.logOutErr(out.String(), ser.String())
		return c.JSON(http.StatusInternalServerError, err)
	}
	bs.logOutErr(out.String(), ser.String())
	return c.JSON(http.StatusOK, "OK")
}

func (bs *barService) status(c echo.Context) error {
	sS := ""
	if !bs.started.IsZero() {
		sS = bs.started.UTC().Format("2006-01-02T15:04:05.000")
	}
	eS := ""
	if !bs.ended.IsZero() {
		eS = bs.ended.UTC().Format("2006-01-02T15:04:05.000")
	}
	status := mStatus{
		Current:       bs.currentStatus,
		LastOperation: bs.lastOperationStatus,
		Started:       sS,
		Ended:         eS,
	}
	return c.JSON(http.StatusOK, status)
}

func (bs *barService) Name() string {
	return "barsvc"
}

func (bs *barService) Start() error {
	bs.Logger().SetLevel(log.INFO)
	e := bs.e
	e.POST("/backup", func(c echo.Context) error {
		return bs.bor(c, false)
	})
	e.POST("/restore", func(c echo.Context) error {
		return bs.bor(c, true)
	})
	e.GET("/status", func(c echo.Context) error {
		return bs.status(c)
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

func (bs *barService) Logger() *log.Logger {
	return bs.logger
}

func newSvc() svcctl.ControllableService {
	bs := &barService{e: echo.New(), logger: log.New("barsvc")}
	bs.e.Use(middleware.Logger())
	bs.configFromEnv()
	return bs
}

func BarSvc() svcctl.ControllableService {
	return newSvc()
}
