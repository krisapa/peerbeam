package receiver

import (
	"github.com/krisapa/peerbeam/utils"
	"github.com/pion/webrtc/v4"
	log "github.com/sirupsen/logrus"
)

func (r *Receiver) SetupReceiverConn() error {
	err := r.Session.SetupPeerConn()
	if err != nil {
		return err
	}
	r.registerHandlers()
	return nil
}

func (r *Receiver) registerHandlers() {
	r.Session.Conn.OnDataChannel(func(ch *webrtc.DataChannel) {
		switch ch.Label() {
		case "data":
			r.Session.DataChHandler(ch)
		default:
			log.Errorln("Unknown channel label:", ch.Label())
			return
		}
	})
}

func (r *Receiver) CreateAnswer(offer string) (string, error) {
	offerSDP, candidates, err := utils.DecodeSDP(offer)
	if err != nil {
		return "", err
	}
	err = r.Session.Conn.SetRemoteDescription(*offerSDP)
	if err != nil {
		return "", err
	}

	initAnswer, err := r.Session.Conn.CreateAnswer(nil)
	if err != nil {
		return "", err
	}
	err = r.Session.Conn.SetLocalDescription(initAnswer)
	if err != nil {
		return "", err
	}
	r.Session.AddCandidates(candidates)
	r.Session.WaitGatherComplete()

	answer := r.Session.Conn.LocalDescription()

	encodedSDP, err := utils.EncodeSDP(answer, r.Session.Candidates)
	if err != nil {
		return "", err
	}

	return encodedSDP, nil
}
