package barsvc

import (
	"context"
	"fmt"
	"github.com/t-beigbeder/otvl_devops_tools/src/go/bar_service/svcctl"
	"os"
	"testing"
)

func TestBarSvc(t *testing.T) {
	go func() {
		err := svcctl.UnderControl(context.Background(), BarSvc(), 10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}()
	//c = http.
}
