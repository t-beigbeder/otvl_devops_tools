package ht3mock

import (
	"errors"
	"fmt"
	"ht3mock/pkg/util"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"sync"
	"time"

	quic "github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type HmClient interface {
	RunTest() error
}

type exexCtx struct {
	runPosition int
	dSize       int
	delay       time.Duration
	nbPending   int
	err         error
}

type runResult struct {
	started, ended time.Time
	run            int
	err            error
}

type hmSingleClient struct {
	p          *hmClient
	index      int
	pat        *WlPattern
	lock       sync.Mutex
	pending    int
	pwg        sync.WaitGroup
	rwg        sync.WaitGroup
	results    map[int]runResult
	done       bool
	terminated bool
}

type hmsResult struct {
	index int
	err   error
}

type hmClient struct {
	url        string
	pats       []*WlPattern
	client     *http.Client
	logger     *slog.Logger
	hmsc       []*hmSingleClient
	done       sync.WaitGroup
	hmsResults chan hmsResult
}

func (hms *hmSingleClient) setDone() {
	hms.lock.Lock()
	defer hms.lock.Unlock()
	hms.done = true
}

func (hms *hmSingleClient) isDone() bool {
	hms.lock.Lock()
	defer hms.lock.Unlock()
	return hms.done
}

func (hms *hmSingleClient) setTerminated() {
	hms.lock.Lock()
	defer hms.lock.Unlock()
	hms.terminated = true
}

func (hms *hmSingleClient) isTerminated() bool {
	hms.lock.Lock()
	defer hms.lock.Unlock()
	return hms.terminated
}

func (hms *hmSingleClient) newRun() {
	hms.lock.Lock()
	hms.pending++
	hms.rwg.Add(1)
	if hms.pending == hms.pat.MaxPending {
		hms.pwg.Add(1)
		hms.lock.Unlock()
		hms.pwg.Wait()
		return
	}
	hms.lock.Unlock()
}

func (hms *hmSingleClient) endRun(rr runResult) {
	hms.lock.Lock()
	hms.pending--
	hms.rwg.Done()
	if hms.pending == hms.pat.MaxPending-1 {
		hms.pwg.Done()
	}
	hms.results[rr.run] = rr
	hms.lock.Unlock()
}

func (hms *hmSingleClient) RunTest() error {
	for run := 0; run < hms.pat.RunNb || hms.pat.RunNb == -1; run++ {
		if hms.isDone() {
			break
		}
		hms.newRun()
		go func(run int) {
			rr := runResult{started: time.Now(), run: run}
			rr.err = hms.run(run)
			rr.ended = time.Now()
			hms.endRun(rr)
		}(run)
	}
	hms.rwg.Wait()
	ks := make([]int, 0, len(hms.results))
	for k := range hms.results {
		ks = append(ks, k)
	}
	sort.Ints(ks)
	for k := range ks {
		if hms.results[k].err != nil {
			rr := hms.results[k]
			hms.p.logger.Error("hmock error", slog.Time("started", rr.started), slog.Time("ended", rr.ended), slog.String("err", rr.err.Error()))
		}
	}
	return nil
}

func (hms *hmSingleClient) run(runPos int) error {
	ctx := exexCtx{runPosition: runPos}
	ctx.dSize = rnd(hms.pat.MinReqDsize, hms.pat.MaxReqDzize)
	ctx.delay = rndDelay(hms.pat.ReqMinDelay, hms.pat.ReqMaxDelay)
	time.Sleep(ctx.delay)
	sDSize := rnd(hms.pat.MinRspDsize, hms.pat.MaxRspDzize)
	sDelay := rnd(hms.pat.RspMinDelay, hms.pat.RspMaxDelay)
	url := fmt.Sprintf("%s?name=%s&pos=%d&dsize=%d&delay=%d", hms.p.url, hms.pat.Name, runPos, sDSize, sDelay)
	bf := util.NewRandGenerator(ctx.dSize, hms.p.logger)
	rsp, err := hms.p.client.Post(url, "application/ octet-stream", bf)
	if err != nil {
		return err
	}
	n, err := io.Copy(io.Discard, rsp.Body)
	if err != nil {
		return err
	}
	hms.p.logger.Info("hmSingleClient.run", slog.Int("sent", ctx.dSize), slog.Int64("received", n), slog.Duration("delay", ctx.delay))
	return nil
}

func newHmSingleClient(index int, p *hmClient, pat *WlPattern) *hmSingleClient {
	h := &hmSingleClient{index: index, p: p, pat: pat}
	pat.initialize()
	h.results = make(map[int]runResult, h.pat.RunNb)
	return h
}

func (hm *hmClient) RunTest() error {
	for _, hms := range hm.hmsc {
		hm.logger.Info("client starting", slog.String("name", hms.pat.Name))
		hm.done.Add(1)

		go func() {
			hms.p.logger.Debug("client start", slog.String("name", hms.pat.Name))
			err := hms.RunTest()
			hms.p.logger.Debug("client done", slog.String("name", hms.pat.Name))
			hms.p.hmsResults <- hmsResult{
				index: hms.index,
				err:   err,
			}
		}()
	}
	var err error
	go func() {
		for {
			select {
			case hmsr := <-hm.hmsResults:
				hms := hm.hmsc[hmsr.index]
				hm.logger.Info("client finished", slog.String("name", hms.pat.Name))
				err = errors.Join(hmsr.err)
				hms.setTerminated()
				allDone := true
				for _, other := range hm.hmsc {
					if !other.isTerminated() && other.pat.RunNb != -1 {
						allDone = false
					}
				}
				if allDone {
					for _, other := range hm.hmsc {
						if !other.isTerminated() && other.pat.RunNb == -1 {
							other.setDone()
						}
					}
				}
				hm.done.Done()
			}
		}
	}()
	hm.done.Wait()
	return err
}

func NewHmClient(config *Config, client *http.Client, logger *slog.Logger) (HmClient, error) {
	hm := &hmClient{client: client, logger: logger}
	hm.url = fmt.Sprintf("https://%s:%d", config.ServerHost, config.Port)
	hm.hmsResults = make(chan hmsResult)
	if len(config.Patterns) == 0 {
		config.Patterns = []WlPattern{{}}
	}
	hm.hmsc = make([]*hmSingleClient, 0)
	for index, pat := range config.Patterns {
		hm.hmsc = append(hm.hmsc, newHmSingleClient(index, hm, &pat))
	}
	return hm, nil
}

func RunClient(config *Config, logger *slog.Logger) error {
	var hc *http.Client
	tc, err := getTlsClientConfig(config)
	if err != nil {
		return err
	}
	if !config.IsHttp2 {
		roundTripper := &http3.RoundTripper{
			TLSClientConfig: tc,
			QUICConfig:      &quic.Config{},
		}
		hc = &http.Client{
			Transport: roundTripper,
		}
	} else {
		hc = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tc,
			},
		}
	}
	c, err := NewHmClient(config, hc, logger)
	if err != nil {
		return err
	}
	if err := c.RunTest(); err != nil {
		return err
	}
	return nil
}
