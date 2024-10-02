package conn

import (
	"context"
	"github.com/pion/webrtc/v4"
	"sync"
	"sync/atomic"
)

type Session struct {
	Conn *webrtc.PeerConnection

	candidateCh     *webrtc.DataChannel
	candidateChOpen atomic.Bool
	CandidateCond   *sync.Cond

	DataCh     *webrtc.DataChannel
	DataChOpen chan struct{}

	Ctx       context.Context
	CtxCancel context.CancelFunc

	MsgCh chan *webrtc.DataChannelMessage
}

func New() *Session {
	ctx, cancel := context.WithCancel(context.Background())
	return &Session{
		Ctx:       ctx,
		CtxCancel: cancel,
		//StateMap:   make(map[webrtc.ICEConnectionState]chan struct{}),
		MsgCh:         make(chan *webrtc.DataChannelMessage, 100),
		DataChOpen:    make(chan struct{}, 10),
		CandidateCond: sync.NewCond(&sync.Mutex{}),
	}
}
