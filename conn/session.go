package conn

import (
	"context"
	"github.com/pion/webrtc/v4"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Session struct {
	Conn *webrtc.PeerConnection

	DataCh *webrtc.DataChannel

	DataChOpen chan struct{}

	Ctx       context.Context
	CtxCancel context.CancelFunc

	MsgCh          chan *webrtc.DataChannelMessage
	Candidates     []webrtc.ICECandidateInit
	CandidatesLock sync.Mutex
	GatherDone     chan struct{}
}

func New() *Session {
	ctx, cancel := context.WithCancel(context.Background())
	return &Session{
		Ctx:        ctx,
		CtxCancel:  cancel,
		DataChOpen: make(chan struct{}, 10),
		GatherDone: make(chan struct{}, 1),
		MsgCh:      make(chan *webrtc.DataChannelMessage, 200),
	}
}

func (c *Session) WaitGatherComplete() {
	<-c.GatherDone
}
func (c *Session) AddCandidates(candidates []webrtc.ICECandidateInit) {
	for i := range candidates {
		err := c.Conn.AddICECandidate(candidates[i])
		if err != nil {
			log.Errorf("add ICECandidate error: %v candidate:%v", err, candidates[i])
			continue
		} else {
			log.Debugf("add ICECandidate success candidate:%v", candidates[i])
		}
	}
}
