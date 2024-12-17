package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/qutils"
	"bufio"
	"fmt"
	"github.com/quic-go/quic-go"
	"golang.org/x/net/context"
	"os"
)

func provision(config *bssms.ProvisionerConfig, stream quic.Stream) error {
	_, err := stream.Write([]byte(bssms.ProvisionerHello))
	if err != nil {
		return err
	}
	rs := bufio.NewReaderSize(stream, bssms.CtrlMsgMaxLn)
	cmd, err := rs.ReadString('\n')
	if err != nil {
		return err
	}
	if cmd != bssms.ProxyHello {
		return fmt.Errorf("invalid protocol command %s", cmd)
	}
	_, err = stream.Write([]byte(bssms.ApplicationClose))
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
	fmt.Fprintf(os.Stderr, "sid %v\n", stream.StreamID())
	defer stream.Close()
	return provision(config, stream)
}
