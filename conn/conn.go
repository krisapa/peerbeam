package conn

import (
	"fmt"
	"github.com/pion/logging"
	"github.com/pion/webrtc/v4"
	log "github.com/sirupsen/logrus"
	"os"
)

func (c *Session) SetupPeerConn() error {
	file, err := os.OpenFile("log/webrtc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	loggerFactory := logging.NewDefaultLoggerFactory()
	loggerFactory.Writer = file
	loggerFactory.DefaultLogLevel = logging.LogLevelTrace

	s := webrtc.SettingEngine{
		LoggerFactory: loggerFactory,
	}
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))
	conn, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers:         iceServers,
		ICETransportPolicy: webrtc.ICETransportPolicyAll,
	})
	//conn, err := webrtc.NewPeerConnection(webrtc.Configuration{
	//	ICEServers: []webrtc.ICEServer{iceServers[0]},
	//})

	if err != nil {
		return err
	}
	c.Conn = conn
	c.monitorState()
	go c.connClose()
	c.sendCandidatesHandler()

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
