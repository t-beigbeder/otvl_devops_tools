package common

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quic-go/quic-go"
	"io"
	"strconv"
	"strings"
)

const (
	CtrlDataMaxLn = 8192
)

func ReadBytesFromStream(sbr *bufio.Reader) ([]byte, error) {
	sln, err := sbr.ReadString('\n')
	if err != nil {
		return nil, err
	}
	ln, err := strconv.ParseInt(sln[0:len(sln)-1], 10, 64)
	if err != nil {
		return nil, err
	}
	if ln > CtrlDataMaxLn {
		return nil, fmt.Errorf("received data length %d > %d", ln, CtrlDataMaxLn)
	}
	bs := make([]byte, ln)
	rln, err := io.ReadFull(sbr, bs)
	if err != nil {
		return nil, err
	}
	if int64(rln) != ln {
		return nil, fmt.Errorf("read data length %d != %d", rln, ln)
	}
	return bs, nil
}

func ReadJsonFromStream(sbr *bufio.Reader, v any) error {
	var (
		js  []byte
		err error
	)
	if js, err = ReadBytesFromStream(sbr); err != nil {
		return err
	}
	if err = json.Unmarshal(js, v); err != nil {
		return err
	}
	return nil
}

func WriteCmdBytesToStream(sbr *bufio.Reader, sw io.Writer, cmd string, payload []byte, ans string) error {
	var (
		err error
	)
	_, err = sw.Write([]byte(cmd))
	if err != nil {
		return err
	}
	if payload != nil {
		_, err = sw.Write([]byte(fmt.Sprintf("%d\n", len(payload))))
		if err != nil {
			return err
		}
		_, err = sw.Write(payload)
		if err != nil {
			return err
		}
	}
	if ans != "" {
		cmd, err := sbr.ReadString('\n')
		if err != nil {
			var ae *quic.ApplicationError
			if !errors.As(err, &ae) || !ae.Remote || ae.ErrorCode != 0 || !strings.Contains(ans, "Bye") {
				return err
			}
			return nil
		}
		if cmd != ans {
			return fmt.Errorf("invalid protocol command %s", cmd)
		}
	}
	return nil
}

func WriteCommandToStream(sbr *bufio.Reader, sw io.Writer, cmd string, payload any, ans string) error {
	var (
		bs  []byte
		err error
	)
	if payload != nil {
		bs, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}
	return WriteCmdBytesToStream(sbr, sw, cmd, bs, ans)
}

func WriteByeCommandToStream(sbr *bufio.Reader, sw io.Writer, cmd string, ans string) (string, error) {
	var (
		cmdr string
		err  error
	)
	_, err = sw.Write([]byte(cmd))
	if err != nil {
		return "", err
	}
	cmdr, err = sbr.ReadString('\n')
	if err != nil {
		var ae *quic.ApplicationError
		if !errors.As(err, &ae) || !ae.Remote || ae.ErrorCode != 0 {
			return "", err
		}
		return "", nil
	}
	if cmdr != ans {
		return "", fmt.Errorf("invalid protocol command %s", cmdr)
	}
	return cmdr, nil
}
