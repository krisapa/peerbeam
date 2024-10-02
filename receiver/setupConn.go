package receiver

import (
	"github.com/pion/webrtc/v4"
	"log/slog"
)

func (r *Receiver) addChHandler() {
	r.Conn.OnDataChannel(func(ch *webrtc.DataChannel) {
		switch ch.Label() {
		case "data":
			r.DataChHandler(ch)
		case "candidate":
			r.CandidateChHandler(ch)
		default:
			slog.Error("Unknown channel label:", ch.Label())
			return
		}
	})
}
