package popup

import (
	"net"
	"sync/atomic"
)

var dialSessionId uint64

func Listen(network, laddr string) (*Server, error) {
	listener, err := net.Listen(network, laddr)
	if err != nil {
		return nil, err
	}
	return NewServer(listener), nil
}


func Dial(network, laddr string) (*Session, error) {
	conn, err := net.Dial(network, laddr)
	if err != nil {
		return nil, err
	}
	id := atomic.AddUint64(&dialSessionId, 1)
	session := NewSession(id, conn)
	return session, nil
}

