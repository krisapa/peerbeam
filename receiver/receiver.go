package receiver

import (
	"github.com/6b70/peerbeam/conn"
)

type Receiver struct {
	Session *conn.Session
}

func New() *Receiver {
	return &Receiver{
		Session: conn.New(),
	}
}
