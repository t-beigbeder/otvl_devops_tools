package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/sleep", func(rw http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			herr(rw, err)
			return
		}
		sd := request.FormValue("delay")
		d, err := strconv.ParseInt(sd, 10, 64)
		if err != nil {
			herr(rw, err)
			return
		}
		start := time.Now().Format("2006-01-02 15:04:05.000")
		time.Sleep(time.Duration(d) * time.Second)
		end := time.Now().Format("2006-01-02 15:04:05.000")
		hmsg(rw, start, end, "ok")
	})
	http.HandleFunc("/healthz", okh)
	http.HandleFunc("/readyz", okh)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Error(fmt.Sprintf("%v", http.ListenAndServe(":8080", nil)))
}

func okh(rw http.ResponseWriter, rq *http.Request) {
	hmsg(rw, "", "", "ok")
}

type MsgResponse struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
	Msg   string `json:"msg"`
}

func hmsg(rw http.ResponseWriter, start, end, msg string) {
	bs, _ := json.Marshal(&MsgResponse{
		Start: start,
		End:   end,
		Msg:   msg,
	})
	_, _ = rw.Write(bs)
}

type ErrResponse struct {
	Error string `json:"error"`
}

func herr(rw http.ResponseWriter, err error) {
	bs, _ := json.Marshal(&ErrResponse{
		Error: fmt.Sprintf("%v", err),
	})
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write(bs)
}
