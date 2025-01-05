package main

import (
	"log"
	"net/http"
	"time"

	socketio "github.com/smart-kf/go-socket.io"
	"github.com/smart-kf/go-socket.io/engineio"
	"github.com/smart-kf/go-socket.io/engineio/transport"
	"github.com/smart-kf/go-socket.io/engineio/transport/websocket"
)

func main() {
	// wss://goim.smartkf.top/sub/?token=helloworld&EIO=3&transport=websocket
	// https://goim.smartkf.top/socket.io/?token=helloworld&EIO=3&transport=websocket
	// Simple client to talk to default-http example
	uri := "http://goim.smartkf.top/socket.io"

	client, err := socketio.NewClient(
		uri, &engineio.Options{
			Transports: []transport.Transport{
				&websocket.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
		},
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

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	client.Emit("notice", "hello")

	time.Sleep(1 * time.Second)
	err = client.Close()
	if err != nil {
		panic(err)
	}
}
