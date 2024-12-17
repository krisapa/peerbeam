package conn

import (
	"fmt"
	"time"
)

func (c *Session) ReceiveMessage(timeout time.Duration) ([]byte, error) {
	select {
	case <-c.Ctx.Done():
		return nil, c.Ctx.Err()
	case <-time.After(timeout):
		return nil, fmt.Errorf("receive timeout")
	case msg := <-c.MsgCh:
		return msg.Data, nil
	}
}
