package main

import (
	"fmt"
	"log"

	socketio "github.com/smart-kf/go-socket.io"
)

func main() {
	// wss://goim.smartkf.top/sub/?token=helloworld&EIO=3&transport=websocket
	client, err := socketio.NewClient(
		"wss://goim.smartkf.top/sub/?token=helloworld&EIO=3&transport=websocket", nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	client.OnConnect(
		func(c socketio.Conn) error {
			fmt.Println("connect")
			fmt.Println(c.ID())
			return nil
		},
	)

	client.OnEvent(
		"msg", func(s socketio.Conn, msg string) string {
			s.SetContext(msg)
			return "recv " + msg
		},
	)

	client.Connect()

	client.Emit("message", "hello world")

}
