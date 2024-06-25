package svcctl

import (
	"context"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ControllableService interface {
	Name() string
	Start() error
	Stop(ctx context.Context) error
	Logger() *log.Logger
}

func UnderControl(ctx context.Context, cs ControllableService, seconds int) error {
	cs.Logger().Infof("starting service %s control", cs.Name())
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancelSvc := context.WithCancel(ctx)
	go func() {
		defer cancelSvc()
		err := cs.Start()
		if err != nil {
			cs.Logger().Error(err)
		}
		cs.Logger().Infof("service %s start finished", cs.Name())
	}()
	<-ctx.Done()
	cs.Logger().Infof("starting service %s shutdown", cs.Name())
	ctx, cancelSD := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancelSD()
	if err := cs.Stop(ctx); err != nil {
		cs.Logger().Error(err)
		return err
	}
	cs.Logger().Infof("service %s has been shutdown, finishing", cs.Name())
	return nil
}
