package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const CtrlDataMaxLn = 8192

func ReadJsonFromStream(streamReader *bufio.Reader, v any) error {
	sln, err := streamReader.ReadString('\n')
	if err != nil {
		return err
	}
	ln, err := strconv.ParseInt(sln[0:len(sln)-1], 10, 64)
	if err != nil {
		return err
	}
	if ln > CtrlDataMaxLn {
		return fmt.Errorf("received data length %d > %d", ln, CtrlDataMaxLn)
	}
	js := make([]byte, ln)
	rln, err := io.ReadFull(streamReader, js)
	if err != nil {
		return err
	}
	if int64(rln) != ln {
		return fmt.Errorf("read data length %d != %d", rln, ln)
	}
	if err = json.Unmarshal(js, v); err != nil {
		return err
	}
	return nil
}

func WriteCommandToStream(stream io.ReadWriter, cmd string, v any) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = stream.Write([]byte(cmd))
	if err != nil {
		return err
	}
	_, err = stream.Write([]byte(fmt.Sprintf("%d\n", len(bs))))
	if err != nil {
		return err
	}
	_, err = stream.Write(bs)
	if err != nil {
		return err
	}
	return nil
}
