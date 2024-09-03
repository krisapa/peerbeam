package receiver

import (
	"fmt"
	"github.com/pion/webrtc/v4"
)

func (r *Receiver) addChHandler() {
	r.Conn.OnDataChannel(func(ch *webrtc.DataChannel) {
		switch ch.Label() {
		case "data":
			r.dataChHandler(ch)
		default:
			fmt.Println("Unknown channel label:", ch.Label())
			return
		}
	})
}

func (r *Receiver) dataChHandler(ch *webrtc.DataChannel) {
	r.DataCh = ch
	ch.OnOpen(func() {
		//fmt.Println("Data Channel Opened")
		r.DataChOpen <- struct{}{}
	})
	r.DataCh.OnClose(func() {
		r.CtxCancel()
	})

	r.DataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
		//fmt.Println("Data Ch Message Received")
		r.MsgCh <- &msg
	})

}
