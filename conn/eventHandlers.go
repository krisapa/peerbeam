package conn

import (
	"encoding/json"
	"github.com/pion/webrtc/v4"
	log "github.com/sirupsen/logrus"
)

func (c *Session) DataChHandler(ch *webrtc.DataChannel) {
	c.DataCh = ch
	ch.OnOpen(func() {
		c.DataChOpen <- struct{}{}
	})
	c.DataCh.OnClose(func() {
		c.CtxCancel()
	})

	c.DataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
		c.MsgCh <- &msg
	})
}

func (c *Session) CandidateChHandler(ch *webrtc.DataChannel) {
	c.candidateCh = ch
	ch.OnOpen(func() {
		c.candidateChOpen.Store(true)
	})

	ch.OnMessage(func(msg webrtc.DataChannelMessage) {
		var candidate webrtc.ICECandidateInit
		err := json.Unmarshal(msg.Data, &candidate)
		if err != nil {
			log.Errorln("Error unmarshalling candidate:", err)
			return
		}
		err = c.Conn.AddICECandidate(candidate)
		if err != nil {
			log.Errorln("Error adding ice candidate:", err)
		}
	})
}

func (c *Session) sendCandidatesHandler() {
	c.Conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		c.CandidateCond.L.Lock()
		c.CandidateCond.Broadcast()
		c.CandidateCond.L.Unlock()

		if !c.candidateChOpen.Load() {
			return
		}

		candidateBytes, err := json.Marshal(candidate.ToJSON())
		if err != nil {
			log.Errorln("Error marshalling candidate:", err)
			return
		}
		err = c.candidateCh.Send(candidateBytes)
		if err != nil {
			log.Errorln("Error sending candidate:", err)
		}
	})
}

func (c *Session) monitorState() {
	c.Conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		switch connectionState {
		case webrtc.ICEConnectionStateFailed:
			c.CtxCancel()
		default:
		}
	})
}
