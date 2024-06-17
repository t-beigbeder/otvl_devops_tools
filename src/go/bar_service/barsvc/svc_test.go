package barsvc

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/t-beigbeder/otvl_devops_tools/src/go/bar_service/svcctl"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestBarSvcBasic(t *testing.T) {
	checkResp := func(resp *http.Response, err error) string {
		if err != nil {
			t.Fatal(err)
		}
		buf := new(strings.Builder)
		io.Copy(buf, resp.Body)
		res := fmt.Sprintf("%d %s", resp.StatusCode, buf.String())
		fmt.Println(res)
		return res
	}
	bs := BarSvc()
	bars, _ := bs.(*barService)
	bars.configure(":3000", []string{"sh", "-c", "echo backup;sleep 1"}, []string{"sh", "-c", "echo restore"})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := svcctl.UnderControl(ctx, bs, 10)
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Millisecond) // server startup
	hc := &http.Client{}
	jrs := checkResp(hc.Post("http://localhost:3000/backup", echo.MIMEApplicationJSON, nil))
	_ = jrs
	time.Sleep(100 * time.Millisecond) // shell command execution
	jrs = checkResp(hc.Get("http://localhost:3000/status"))
	jrs = checkResp(hc.Post("http://localhost:3000/restore", echo.MIMEApplicationJSON, nil))
	jrs = checkResp(hc.Get("http://localhost:3000/status"))
	jrs = checkResp(hc.Get("http://localhost:3000/healthz"))
	time.Sleep(1100 * time.Millisecond)
	jrs = checkResp(hc.Post("http://localhost:3000/restore", echo.MIMEApplicationJSON, nil))
	jrs = checkResp(hc.Get("http://localhost:3000/status"))
	time.Sleep(100 * time.Millisecond) // shell command execution
	jrs = checkResp(hc.Get("http://localhost:3000/status"))
	cancel()
	time.Sleep(100 * time.Millisecond) // server shutdown
}

func TestBarSvcErr(t *testing.T) {
	checkResp := func(resp *http.Response, err error) string {
		if err != nil {
			t.Fatal(err)
		}
		buf := new(strings.Builder)
		io.Copy(buf, resp.Body)
		res := fmt.Sprintf("%d %s", resp.StatusCode, buf.String())
		fmt.Println(res)
		return res
	}
	bs := BarSvc()
	bars, _ := bs.(*barService)
	bars.configure(":3000", []string{"sh", "-c", "ls /tmp/zgdjsgdfqgdqsjg"}, []string{"sh", "-c", "echo restore"})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := svcctl.UnderControl(ctx, bs, 10)
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Millisecond) // server startup
	hc := &http.Client{}
	jrs := checkResp(hc.Post("http://localhost:3000/backup", echo.MIMEApplicationJSON, nil))
	_ = jrs
	time.Sleep(100 * time.Millisecond) // shell command execution
	jrs = checkResp(hc.Get("http://localhost:3000/status"))
	cancel()
	time.Sleep(100 * time.Millisecond) // server shutdown
}
