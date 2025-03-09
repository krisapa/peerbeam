package conn

import (
	"fmt"
	"github.com/pion/webrtc/v4"
	log "github.com/sirupsen/logrus"
)

func (c *Session) SetupPeerConn() error {
	conn, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{iceServers[0]},
	})
	if err != nil {
		return err
	}
	c.Conn = conn
	c.monitorConnectState()
	c.monitorGatherState()
	go c.connClose()
	c.CandidatesHandler()
	return nil
}

func (c *Session) connClose() {
	<-c.Ctx.Done()

	if c.DataCh != nil {
		err := c.DataCh.GracefulClose()
		if err != nil {
			log.Error("Error closing control channel:", err)
		}
	}
	if c.Conn == nil {
		err := c.Conn.GracefulClose()
		if err != nil {
			log.Errorln("Error closing peer connection:", err)
		}
	}
}

func (c *Session) PeerInfoStr() (string, error) {
	selectedPair, err := c.DataCh.Transport().Transport().ICETransport().GetSelectedCandidatePair()
	if err != nil {
		return "", err
	}
	remote := selectedPair.Remote
	return fmt.Sprintf("%s %s:%d", remote.Protocol, remote.Address, remote.Port), nil
}
