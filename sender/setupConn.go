package sender

func (s *Sender) setupDataCh() error {
	ch, err := s.Conn.CreateDataChannel("data", nil)
	if err != nil {
		return err
	}
	s.DataChHandler(ch)
	//s.DataCh = channel
	//s.DataCh.OnOpen(func() {
	//	s.DataChOpen <- struct{}{}
	//})
	//s.DataCh.OnClose(func() {
	//	s.CtxCancel()
	//})
	//
	//s.DataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
	//	s.MsgCh <- &msg
	//})
	return nil
}
