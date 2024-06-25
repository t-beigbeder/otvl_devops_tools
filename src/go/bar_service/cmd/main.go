package main

import (
	"context"
	"fmt"
	"github.com/t-beigbeder/otvl_devops_tools/src/go/bar_service/barsvc"
	"github.com/t-beigbeder/otvl_devops_tools/src/go/bar_service/svcctl"
	"os"
)

func main() {
	err := svcctl.UnderControl(context.Background(), barsvc.BarSvc(), 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
