package popup

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	listener     net.Listener
	sessions     map[uint64]*Session
	maxSessionId uint64
}

func NewServer(listener net.Listener) *Server {
	return &Server{
		listener:listener,
		maxSessionId: 0,
		sessions:make(map[uint64]*Session),
	}
}

//处理客户端连接
func (server *Server) EventLoop(handler func(*Session)) {
	for {
		session, err := server.Accept()
		if err != nil {
			fmt.Println("error when accept new connection")
			break
		}
		go handler(session)
	}
	server.Stop()
}

//接受客户端连接
func (server *Server) Accept() (*Session, error) {
	conn, err := server.listener.Accept()
	if err != nil {
		return nil, err
	}
	session := server.newSession(conn)
	return session, nil
}

//停止服务，在发生异常时，或者主动关闭server时
func (server *Server) Stop() {
	server.listener.Close()//关闭监听
	server.closeSessions()//关闭所有会话
}

func (server *Server) newSession(conn net.Conn) *Session {
	id := atomic.AddUint64(&server.maxSessionId, 1)
	session := NewSession(id, conn)
	return session
}

//关闭所有对话
func (server *Server) closeSessions() {
	for _, session := range server.sessions {
		session.Close()
	}
}
