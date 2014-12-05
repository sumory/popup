package popup

import (
	"net"
	"sync/atomic"
)

var dialSessionId uint64

func Listen(n, laddr string) (*Server, error) {
	listener, err := net.Listen(n, laddr)
	if err != nil {
		return nil, err
	}
	return NewServer(listener), nil
}


func Dial(n, laddr string) (*Session, error) {
	conn, err := net.Dial(n, laddr)
	if err != nil {
		return nil, err
	}
	id := atomic.AddUint64(&dialSessionId, 1)
	session := NewSession(id, conn)
	return session, nil
}

