package main

import "github.com/sumory/popup"
import "fmt"
import "time"

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func sessionHandler(session *popup.Session) {
	fmt.Println("session start")
	session.ReadLoop(readHandler)
	fmt.Println("session closed")
}

func readHandler(session *popup.Session, msg []byte) {
	fmt.Println(Now(),session.GetId(), string(msg))
	if string(msg) == "file_finish:" {
		session.Send([]byte("文件传输完毕"))
	}else {
		session.Send([]byte(string(msg)))
	}

}


func main() {
	server, _ := popup.Listen("tcp", "0.0.0.0:8080")
	server.EventLoop(sessionHandler)
}
