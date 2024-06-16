package barsvc

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/t-beigbeder/otvl_devops_tools/src/go/bar_service/svcctl"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestBarSvc(t *testing.T) {
	bs := BarSvc()
	bars, _ := bs.(*barService)
	bars.configure(":3000", "sh -c 'echo backup'", "sh -c 'echo restore'")
	go func() {
		err := svcctl.UnderControl(context.Background(), bs, 10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}()
	time.Sleep(2 * time.Second)
	hc := &http.Client{}
	hc.Post("http://localhost:3000/exit", echo.MIMEApplicationJSON, nil)

}
