package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"ht3mock/pkg/ht3mock"
)

func main() {
	configFile := flag.String("c", "", "path of a config file with workload patterns")
	port := flag.Int("p", 9443, "port number defaults to 9443")
	listen := flag.String("lh", "0.0.0.0", "server listen address, defaults to 0.0.0.0")
	server := flag.String("sh", "localhost", "server address for client, defaults to localhost")
	isH2 := flag.Bool("h2", false, "run in HTTPS/2 and not HTTP/3, defaults to false")
	cert := flag.String("cf", "/tmp/ht3mock.cert", "server certificate file, defaults to /tmp/ht3mock.cert")
	key := flag.String("kf", "/tmp/ht3mock.key", "server key file, defaults to /tmp/ht3mock.key")
	ccert := flag.String("ccf", "/tmp/ht3mock.cert", "client certificate file, optional")
	unsafe := flag.Bool("ut", false, "disable client TLS certificate checking, defaults to false")
	isServer := flag.Bool("svr", false, "run the server only, defaults to false")
	isClient := flag.Bool("cli", false, "run the client only, defaults to false")
	isDebug := flag.Bool("dbg", false, "logs debug")
	flag.Parse()
	pat := &ht3mock.WlPattern{}
	pats := make([]ht3mock.WlPattern, 0)
	_ = pats
	if configFile != nil && *configFile != "" {
		cd, err := os.ReadFile(*configFile)
		if err != nil {
			exErr(err)
		}
		if err := json.Unmarshal(cd, pat); err != nil {
			if err := json.Unmarshal(cd, &pats); err != nil {
				exErr(err)
			}
		} else {
			pats = append(pats, *pat)
		}
	}
	sll := slog.LevelInfo
	if *isDebug {
		sll = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: sll}))
	logger.Info("main: started")
	marshal, err := json.Marshal(pats)
	if err != nil {
		exErr(err)
	}
	logger.Debug("main: we are in debug mode", slog.String("config", string(marshal)))
	config := &ht3mock.Config{
		Port:       *port,
		ListenHost: *listen,
		ServerHost: *server,
		IsHttp2:    *isH2,
		CertFile:   *cert,
		CCertFile:  *ccert,
		UnsafeTls:  *unsafe,
		KeyFile:    *key,
		Patterns:   pats,
	}
	started := time.Now()
	done := make(chan interface{})
	if !*isClient {
		err = ht3mock.RunServer(config, logger)
		if err != nil {
			exErr(err)
		}
		if *isServer {
			<-done
		}
	}
	if !*isServer {
		if err := ht3mock.RunClient(config, logger); err != nil {
			close(done)
			exErr(err)
		}
	}
	close(done)
	logger.Info("took", slog.Duration("duration", time.Since(started)))
}

func exErr(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
	os.Exit(-1)
}
