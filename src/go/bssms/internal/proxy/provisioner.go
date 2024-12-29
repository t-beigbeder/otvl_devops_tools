package proxy

import (
	"bssms/internal/bssms"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func handlePrCmd(rs *bufio.Reader, cmd string) error {
	sln, err := rs.ReadString('\n')
	if err != nil {
		return err
	}
	ln, err := strconv.ParseInt(sln[0:len(sln)-1], 10, 64)
	if err != nil {
		return err
	}
	if ln > bssms.CtrlDataMaxLn {
		return fmt.Errorf("received data length %d > %d", ln, bssms.CtrlDataMaxLn)
	}
	js := make([]byte, ln)
	rln, err := io.ReadFull(rs, js)
	if err != nil {
		return err
	}
	if int64(rln) != ln {
		return fmt.Errorf("read data length %d != %d", rln, ln)
	}
	in := bssms.Installable{}
	if err = json.Unmarshal(js, &in); err != nil {
		return err
	}
	getLogger().Debug("handlePrCmd", "in", in)
	return nil
}
