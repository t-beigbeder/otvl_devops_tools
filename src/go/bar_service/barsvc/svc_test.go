package barsvc

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/t-beigbeder/otvl_devops_tools/src/go/bar_service/svcctl"
	"net/http"
	"os"
	"testing"
)

func TestBarSvc(t *testing.T) {
	bs := BarSvc()
	bars, _ := bs.(*barService)
	bars.configure(":3000", []string{"sh", "-c", "echo backup"}, []string{"sh", "-c", "echo restore"})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := svcctl.UnderControl(ctx, bs, 10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}()
	hc := &http.Client{}
	hc.Post("http://localhost:3000/backup", echo.MIMEApplicationJSON, nil)
	hc.Post("http://localhost:3000/restore", echo.MIMEApplicationJSON, nil)
	hc.Get("http://localhost:3000/status")
	hc.Get("http://localhost:3000/healthz")
	cancel()
}
