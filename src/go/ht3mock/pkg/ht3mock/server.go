package ht3mock

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"ht3mock/pkg/util"

	"github.com/quic-go/quic-go/http3"
)

func RunServer(config *Config, logger *slog.Logger) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		bytesRead, err := io.Copy(io.Discard, request.Body)
		if err != nil || (request.ContentLength != -1 && bytesRead != request.ContentLength) {
			writer.WriteHeader(http.StatusInternalServerError)
			logger.Debug("HandleFunc: copy from", slog.Int64("bytesRead", bytesRead), util.ErrAttr(err))
			writer.Write(
				[]byte(fmt.Sprintf("bytesRead %d err %v ContentLength %d",
					bytesRead, err, request.ContentLength)))
			return
		}
		patName := request.FormValue("name")
		pos := request.FormValue("pos")
		dSize, _ := strconv.ParseInt(request.FormValue("dsize"), 10, 64)
		msDelay, _ := strconv.ParseInt(request.FormValue("delay"), 10, 64)
		delay := time.Duration(msDelay) * time.Millisecond
		logger.Debug("HandleFunc: sleep", slog.String("name", patName), slog.String("pos", pos), slog.Int64("bytesRead", bytesRead), slog.Int64("dSize", dSize), slog.Duration("delay", delay))
		time.Sleep(delay)
		writer.WriteHeader(http.StatusOK)
		bf := util.NewRandGenerator(int(dSize), logger)
		n, err := io.Copy(writer, bf)
		if err != nil {
			logger.Debug("HandleFunc: copy to", slog.Int64("n", n), util.ErrAttr(err))
			return
		}
		logger.Debug("HandleFunc: done", slog.Int("dSize", int(dSize)), slog.Duration("delay", delay))
		return
	})
	addr := fmt.Sprintf("%s:%d", config.ListenHost, config.Port)
	go func() {
		var err error
		loggingMiddleware := util.LoggingMiddleware(logger)
		loggedMux := loggingMiddleware(mux)
		if !config.IsHttp2 {
			err = http3.ListenAndServeQUIC(addr, config.CertFile, config.KeyFile, loggedMux)
		} else {
			err = http.ListenAndServeTLS(addr, config.CertFile, config.KeyFile, loggedMux)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
	}()
	time.Sleep(200 * time.Millisecond)
	return nil
}
