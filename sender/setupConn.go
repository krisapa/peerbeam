package sender

import (
	"github.com/pion/webrtc/v4"
)

func (s *Sender) setupDataCh() error {
	channel, err := s.Conn.CreateDataChannel("data", nil)
	if err != nil {
		return err
	}
	s.DataCh = channel
	s.DataCh.OnOpen(func() {
		s.DataChOpen <- struct{}{}
	})
	s.DataCh.OnClose(func() {
		s.CtxCancel()
	})

	s.DataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
		s.MsgCh <- &msg
	})

	return nil
}
