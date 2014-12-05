package popup

import "net"

type Session struct {
	id       uint64
	conn     net.Conn
	sendChan chan []byte
}

//创建新会话
func NewSession(id uint64, conn net.Conn) *Session {
	session := &Session{
		id:id,
		conn:conn,
		sendChan: make(chan []byte),
	}

	go session.sendLoop()
	return session
}

func (session *Session) GetId() uint64{
	return session.id
}

//关闭会话
func (session *Session) Close() {
	session.conn.Close()
}

//读取数据.
func (session *Session) Read() ([]byte, error) {
	buf := make([]byte, 1024)
	n,err:=session.conn.Read(buf)
	return buf[:n], err
}

func (session *Session) ReadLoop(handler func([]byte)) {
	for {
		var buf []byte
		var err error
		if buf, err = session.Read(); err != nil {
			session.Close()
			break
		}
		handler(buf)
	}
}

func (session *Session) Send(message []byte) error {
	_, err := session.conn.Write(message)
	return err
}


func (session *Session) sendLoop() {
	for {
		message := <-session.sendChan
		if err := session.Send(message); err != nil {
			session.Close()
			return
		}
	}
}
