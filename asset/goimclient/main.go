package main

import (
	"fmt"
	"log"
	"time"

	socketio "github.com/smart-kf/go-socket.io"
)

func main() {
	// wss://goim.smartkf.top/sub/?token=helloworld&EIO=3&transport=websocket
	// https://goim.smartkf.top/socket.io/?token=helloworld&EIO=3&transport=websocket
	// Simple client to talk to default-http example
	uri := "http://goim.smartkf.top/?token=helloworld&platform=kf&transport=websocket"

	client, err := socketio.NewClient(
		uri, nil,
	)
	if err != nil {
		panic(err)
	}

	// Handle an incoming event
	client.OnEvent(
		"reply", func(s socketio.Conn, msg string) {
			log.Println("Receive Message /reply: ", "reply", msg)
		},
	)

	client.OnConnect(
		func(conn socketio.Conn) error {
			fmt.Println(conn.ID())
			return nil
		},
	)

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	client.Emit("notice", "hello")

	time.Sleep(100 * time.Second)
	err = client.Close()
	if err != nil {
		panic(err)
	}
}
