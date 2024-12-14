package main

import (
	"bssms/internal/tlsutils"
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

func main() {

	cert, err := tlsutils.SelfSigned("localhost")
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{*cert}}
	srv := &http.Server{
		Addr:         "0.0.0.0:9443",
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
