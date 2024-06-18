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

type barOperation struct {
	started   time.Time
	operation string
	status    string
	isDone    bool
	hasError  bool
	ended     time.Time
}

type barService struct {
	address         string
	backup, restore []string
	sync            sync.Mutex
	e               *echo.Echo
	logger          *log.Logger
	ongoing         bool
	current         barOperation
	previous        barOperation
}

type mOperationStatus struct {
	Started   string `json:"started,omitempty"`
	Operation string `json:"operation,omitempty"`
	Status    string `json:"status,omitempty"`
	IsDone    bool   `json:"isDone,omitempty"`
	HasError  bool   `json:"hasError,omitempty"`
	Ended     string `json:"ended,omitempty"`
}

type mStatus struct {
	Current  mOperationStatus `json:"current,omitempty"`
	Previous mOperationStatus `json:"previous,omitempty"`
}

type mMsg struct {
	Err string `json:"error,omitempty"`
	Msg string `json:"msg,omitempty"`
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
	ongoing := bs.ongoing
	if !ongoing {
		bs.ongoing = true
		bs.previous = bs.current
		bs.current = barOperation{}
		bs.current.operation = strings.Join(args, " ")
		bs.current.started = time.Now()
	}
	bs.sync.Unlock()
	if ongoing {
		return c.JSON(http.StatusOK, &mMsg{
			Err: fmt.Sprintf("operation in progress: %s", bs.current.operation),
		})
	}
	go func() {
		cmd := exec.Command(args[0], args[1:]...)
		var out, ser strings.Builder
		cmd.Stdout = &out
		cmd.Stderr = &ser
		if err := cmd.Run(); err != nil {
			bs.logOutErr(out.String(), ser.String())
			bs.sync.Lock()
			defer bs.sync.Unlock()
			bs.current.ended = time.Now()
			bs.current.status = fmt.Sprintf("error %v", err)
			bs.current.hasError = true
			bs.current.isDone = true
			bs.ongoing = false
			return
		}
		bs.logOutErr(out.String(), ser.String())
		bs.sync.Lock()
		defer bs.sync.Unlock()
		bs.current.ended = time.Now()
		bs.current.status = ""
		bs.current.isDone = true
		bs.ongoing = false
	}()
	return c.JSON(http.StatusOK, &mMsg{
		Msg: fmt.Sprintf("launched: %s", bs.current.operation),
	})
}

func (bs *barService) status(c echo.Context) error {
	sop := func(op barOperation) mOperationStatus {
		sS := ""
		if !op.started.IsZero() {
			sS = op.started.UTC().Format("2006-01-02T15:04:05.000")
		}
		eS := ""
		if !op.ended.IsZero() {
			eS = op.ended.UTC().Format("2006-01-02T15:04:05.000")
		}
		return mOperationStatus{
			Started:   sS,
			Operation: op.operation,
			Status:    op.status,
			IsDone:    op.isDone,
			HasError:  op.hasError,
			Ended:     eS,
		}
	}
	status := mStatus{
		Current:  sop(bs.current),
		Previous: sop(bs.previous),
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
	bs := &barService{
		e:        echo.New(),
		logger:   log.New("barsvc"),
		current:  barOperation{},
		previous: barOperation{},
	}
	bs.e.Use(middleware.Logger())
	bs.configFromEnv()
	return bs
}

func BarSvc() svcctl.ControllableService {
	return newSvc()
}
