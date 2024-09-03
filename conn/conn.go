package conn

import (
	"context"
	"fmt"
	"github.com/pion/webrtc/v4"
)

func (c *Session) AddRemote(remoteSDP *webrtc.SessionDescription) error {
	err := c.Conn.SetRemoteDescription(*remoteSDP)
	if err != nil {
		return err
	}
	return nil
}

func (c *Session) SetupPeerConn() error {
	conn, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: iceServers,
	})
	if err != nil {
		return err
	}
	c.Conn = conn
	ctx, cancel := context.WithCancel(context.Background())
	c.Ctx = ctx
	c.CtxCancel = cancel
	c.monitorState()
	go c.connClose()

	return nil
}

func (c *Session) connClose() {
	<-c.Ctx.Done()

	if c.DataCh != nil {
		err := c.DataCh.GracefulClose()
		if err != nil {
			fmt.Println("Error closing control channel:", err)
		}
	}
	if c.Conn == nil {
		err := c.Conn.GracefulClose()
		if err != nil {
			fmt.Println("Error closing peer connection:", err)
		}
	}
}

func (c *Session) monitorState() {
	c.Conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		switch connectionState {
		case webrtc.ICEConnectionStateFailed:
			c.CtxCancel()
		default:
			if ch, ok := c.StateMap[connectionState]; ok {
				close(ch)
			}
		}
	})
}
