package network

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/alibaba/ioc-golang/autowire/singleton"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	websocket2 "goim3/application/websocket"
	"goim3/config"
	"goim3/pkg/utils"
)

type WebsocketServer struct {
	socketServer *socketio.Server

	connMutex sync.Mutex
	conns     map[string]socketio.Conn
}

func CreateWebsocketServer() *WebsocketServer {
	wsServer := &WebsocketServer{
		conns: make(map[string]socketio.Conn),
	}

	idGenerator, err := singleton.GetImpl("IDGenerator", nil)
	if err != nil {
		panic(err)
	}

	socketIOServer := socketio.NewServer(
		&engineio.Options{
			PingTimeout:  config.Config.SocketIO.PingTimeout.Duration(),
			PingInterval: config.Config.SocketIO.PingInterval.Duration(),
			Transports: []transport.Transport{
				&websocket.Transport{
					ReadBufferSize:   config.Config.SocketIO.ReadBufferSize,
					WriteBufferSize:  config.Config.SocketIO.WriteBufferSize,
					HandshakeTimeout: config.Config.SocketIO.ReadTimeout.Duration(),
				},
			},
			SessionIDGenerator: idGenerator.(*utils.IDGenerator),
			RequestChecker:     wsServer.onRequestCheck,
			ConnInitor:         wsServer.onConnectionInit,
		},
	)

	wsServer.socketServer = socketIOServer
	wsServer.RegisterEvents()

	go func() {
		if err := wsServer.socketServer.Serve(); err != nil {
			panic(err)
		}
	}()

	return wsServer
}

func (s *WebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.socketServer.ServeHTTP(w, r)
}

// onRequestCheck 建立连接的验证代码
func (s *WebsocketServer) onRequestCheck(r *http.Request) (http.Header, error) {
	fmt.Println("onrequest")
	return http.Header{}, nil
}

// onConnectionInit 验证通过以后得初始化操作.
func (s *WebsocketServer) onConnectionInit(r *http.Request, conn engineio.Conn) {
	return
}

func (s *WebsocketServer) RegisterEvents() {
	s.socketServer.OnConnect(
		"/", func(c socketio.Conn) error {
			s.connMutex.Lock()
			s.conns[c.ID()] = c
			s.connMutex.Unlock()
			c.SetContext("")
			return nil
		},
	)
	s.socketServer.OnEvent("/", "message", s.onMessage)
	s.socketServer.OnDisconnect(
		"/", func(conn socketio.Conn, msg string) {
			s.connMutex.Lock()
			delete(s.conns, conn.ID())
			s.connMutex.Unlock()
			fmt.Println(conn.ID(), "closed", msg)
		},
	)
	s.socketServer.OnError(
		"/", func(c socketio.Conn, err error) {
			s.connMutex.Lock()
			delete(s.conns, c.ID())
			s.connMutex.Unlock()
			fmt.Println(c.ID(), "error", err)
		},
	)
}

func (s *WebsocketServer) onMessage(conn socketio.Conn, msg string) {
	var app websocket2.WebsocketApplication
	app.OnMessage(context.Background(), msg)
}

type PushMessageRequest struct {
	SessionId string `json:"sessionId"`
	Event     string `json:"event"`
	Data      string `json:"data"`
}

func (s *WebsocketServer) Push(ctx *gin.Context) {
	var req PushMessageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(200, gin.H{"err": err})
		return
	}
	s.connMutex.Lock()
	conn, ok := s.conns[req.SessionId]
	s.connMutex.Unlock()
	if !ok {
		return
	}
	conn.Emit(req.Event, req.Data)
}
