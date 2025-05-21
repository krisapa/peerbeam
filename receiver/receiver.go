package receiver

import (
	"github.com/krisapa/peerbeam/conn"
	"time"
)

const DefaultTimeout = 5 * time.Second

type Receiver struct {
	Session *conn.Session
}

func New() *Receiver {
	return &Receiver{
		Session: conn.New(),
	}
}
