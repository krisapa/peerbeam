package conn

import (
	"context"
	"github.com/pion/webrtc/v4"
)

type Session struct {
	Conn       *webrtc.PeerConnection
	DataCh     *webrtc.DataChannel
	DataChOpen chan struct{}

	Ctx       context.Context
	CtxCancel context.CancelFunc

	MsgCh    chan *webrtc.DataChannelMessage
	StateMap map[webrtc.ICEConnectionState]chan struct{}
}

func New() *Session {
	ctx, cancel := context.WithCancel(context.Background())
	return &Session{
		Ctx:        ctx,
		CtxCancel:  cancel,
		StateMap:   make(map[webrtc.ICEConnectionState]chan struct{}),
		MsgCh:      make(chan *webrtc.DataChannelMessage, 100),
		DataChOpen: make(chan struct{}),
	}
}
