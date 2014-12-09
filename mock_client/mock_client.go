package main


import "github.com/sumory/popup"
import "fmt"
import "time"

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}


func main(){
	client, err := popup.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	go client.ReadLoop(readHandler)

	for {
		var input string
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			break
		}
		client.Send([]byte(input))
	}

	client.Close()
}

func readHandler(client *popup.Session, msg []byte) {
	fmt.Println(Now(), "client:",client.GetId(), string(msg))
}
