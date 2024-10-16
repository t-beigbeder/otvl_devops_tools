package htmok

const defaultName = "default"
const defaultDelay = 1
const defaultReqDSize = 100000
const defaultRspDSize = 10000
const defaultMaxPending = 1000
const defaultRunNb = 10000

type WlPattern struct {
	Name        string `json:"name"`
	ReqMinDelay int    `json:"reqMinDelay"`
	ReqMaxDelay int    `json:"reqMaxDelay"`
	RspMinDelay int    `json:"rspMinDelay"`
	RspMaxDelay int    `json:"rspMaxDelay"`
	MinReqDsize int    `json:"minReqDsize"`
	MaxReqDzize int    `json:"maxReqDzize"`
	MinRspDsize int    `json:"minRspDsize"`
	MaxRspDzize int    `json:"maxRspDzize"`
	MaxPending  int    `json:"maxPending"`
	RunNb       int    `json:"runNb"`
}

func (p *WlPattern) initialize() {
	if p.Name == "" {
		p.Name = defaultName
	}
	if p.ReqMinDelay == 0 {
		p.ReqMinDelay = defaultDelay
	}
	if p.ReqMaxDelay == 0 {
		p.ReqMaxDelay = 2 * p.ReqMinDelay
	}
	if p.MinReqDsize == 0 {
		p.MinReqDsize = defaultReqDSize
	}
	if p.MaxReqDzize == 0 {
		p.MaxReqDzize = 2 * p.MinReqDsize
	}
	if p.MinRspDsize == 0 {
		p.MinRspDsize = defaultRspDSize
	}
	if p.MaxRspDzize == 0 {
		p.MaxRspDzize = 2 * p.MinRspDsize
	}
	if p.RspMinDelay == 0 {
		p.RspMinDelay = defaultDelay
	}
	if p.RspMaxDelay == 0 {
		p.RspMaxDelay = 2 * p.RspMinDelay
	}
	if p.MaxPending == 0 {
		p.MaxPending = defaultMaxPending
	}
	if p.RunNb == 0 {
		p.RunNb = defaultRunNb
	}
}

type Config struct {
	Port       int
	ListenHost string
	ServerHost string
	IsHttp2    bool
	CertFile   string
	KeyFile    string
	CCertFile  string
	UnsafeTls  bool
	Patterns   []WlPattern
}
