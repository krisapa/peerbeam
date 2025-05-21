package sender

import (
	"github.com/krisapa/peerbeam/conn"
)

type Sender struct {
	Session *conn.Session
}

func New() *Sender {
	return &Sender{
		Session: conn.New(),
	}
}
