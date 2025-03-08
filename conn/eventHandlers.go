package conn

import (
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

//func (c *Session) CandidateChHandler(ch *webrtc.DataChannel) {
//	c.candidateCh = ch
//	ch.OnOpen(func() {
//		c.candidateChOpen.Store(true)
//	})
//
//	ch.OnMessage(func(msg webrtc.DataChannelMessage) {
//		var candidate webrtc.ICECandidateInit
//		err := json.Unmarshal(msg.Data, &candidate)
//		if err != nil {
//			log.Errorln("Error unmarshalling candidate:", err)
//			return
//		}
//		err = c.Conn.AddICECandidate(candidate)
//		if err != nil {
//			log.Errorln("Error adding ice candidate:", err)
//		}
//	})
//}

func (c *Session) CandidatesHandler() {
	c.Conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			log.Info("find ICECandidate is nil")
			return
		}
		log.Infof("find a new ICECandidate: %+v", candidate.String())
		c.CandidatesLock.Lock()
		defer c.CandidatesLock.Unlock()
		c.Candidates = append(c.Candidates, candidate.ToJSON())
	})
}

func (c *Session) monitorConnectState() {
	c.Conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Info("connection state:", connectionState)
		switch connectionState {
		case webrtc.ICEConnectionStateFailed:
			c.CtxCancel()
		default:
		}
	})
}
func (c *Session) monitorGatherState() {
	c.Conn.OnICEGatheringStateChange(func(state webrtc.ICEGatheringState) {
		log.Info("gathering state:", state.String())
		if state == webrtc.ICEGatheringStateComplete {
			close(c.GatherDone)
		}
	})

}
