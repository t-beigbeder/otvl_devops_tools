package svcctl

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"testing"
	"time"
)

type cableForever struct {
	name   string
	logger echo.Logger
}

func (c *cableForever) Name() string {
	return c.name
}

func (c *cableForever) Start() error {
	c.logger.Debugf("Starting %s", c.name)
	time.Sleep(100 * time.Millisecond)
	c.logger.Debugf("Done %s", c.name)
	return nil
}

func (c *cableForever) Stop(ctx context.Context) error {
	c.logger.Debugf("Stopping %s", c.name)
	return nil
}

func (c *cableForever) Logger() echo.Logger {
	if c.logger == nil {
		c.logger = log.New("bar_service")
	}
	return c.logger
}

func TestControllable(t *testing.T) {
	c := &cableForever{name: "cableForever"}
	c.Logger().SetLevel(log.DEBUG)
	err := UnderControl(context.Background(), c, 1)
	if err != nil {
		t.Fatal(err)
	}
}
