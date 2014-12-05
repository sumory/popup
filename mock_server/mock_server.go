package main

import "github.com/sumory/popup"
import "fmt"
import "time"

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func main(){
	server, _ := popup.Listen("tcp", "0.0.0.0:8080")
	server.EventLoop(func(session *popup.Session) {
		fmt.Println("session start")

		session.ReadLoop(func(msg []byte) {
			fmt.Println(Now(),session.GetId(), string(msg))
			session.Send([]byte("收到："+string(msg)))
		})

		fmt.Println("session closed")
	})
}
