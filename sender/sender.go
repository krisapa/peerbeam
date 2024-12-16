package sender

import (
	"github.com/6b70/peerbeam/conn"
)

type Sender struct {
	Session *conn.Session
}

func New() *Sender {
	return &Sender{
		Session: conn.New(),
	}
}
