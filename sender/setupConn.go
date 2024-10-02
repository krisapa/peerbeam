package sender

func (s *Sender) setupChs() error {
	ch, err := s.Conn.CreateDataChannel("data", nil)
	if err != nil {
		return err
	}
	s.DataChHandler(ch)

	ch, err = s.Conn.CreateDataChannel("candidate", nil)
	if err != nil {
		return err
	}
	s.CandidateChHandler(ch)

	return nil
}
