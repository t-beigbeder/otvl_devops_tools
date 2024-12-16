package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/qutils"
	"fmt"
	"github.com/quic-go/quic-go"
	"golang.org/x/net/context"
	"io"
	"os"
)

func provision(config *bssms.ProvisionerConfig, stream quic.Stream) error {
	_, err := stream.Write([]byte(bssms.ProvisionerHello))
	if err != nil {
		return err
	}
	buf := make([]byte, 128)
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "provision received %v\n", buf)
	_, err = stream.Write([]byte(bssms.MsgClose))
	if err != nil {
		return err
	}
	return nil
}

func Run(config *bssms.ProvisionerConfig) error {
	conn, err := qutils.GetQuicConn(config.ProxyAddress, bssms.BssmsAlpn)
	if err != nil {
		return err
	}
	defer conn.CloseWithError(0, "")
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()
	return provision(config, stream)
}
