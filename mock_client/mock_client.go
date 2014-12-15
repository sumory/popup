package main

import (
	"github.com/sumory/popup"
	"time"
	"strings"
	"os"
	"bufio"
	"io"
	"fmt"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

var canInput chan bool

func main() {
	client, err := popup.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	go client.ReadLoop(readHandler)
	canInput = make(chan bool)
	for {
		var input string
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			break
		}

		command := string(input)
		if strings.HasPrefix(command, "file:") {
			fileName := strings.TrimLeft(command, "file:")
			go func() {
				f, err := os.OpenFile(fileName, os.O_RDONLY, 0660)
				defer f.Close()
				if err != nil {
					panic(err)
				}
				var n int64
				if fi, err := f.Stat(); err == nil {
					if size := fi.Size(); size < 1e9 {
						n = size
					}else {
						//should return error
					}
				}
				fmt.Println("发送文件大小：", n)
				reader := bufio.NewReader(f)
				for {
					buf := make([]byte, 30)
					m, err := reader.Read(buf)
					if err != nil && err != io.EOF {panic(err)}
					if 0 == m {
						break
					} else {
						client.Send(buf[:m])
					}
				}
				canInput <- true
			}()

			<-canInput
			client.Send([]byte("file_finish:"))
		}else {
			client.Send([]byte(input))
		}



	}

	client.Close()
}

func readHandler(client *popup.Session, msg []byte) {
	fmt.Println(Now(), "server已收到:", string(msg))
}
