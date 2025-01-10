package network

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/alibaba/ioc-golang/autowire/singleton"
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	socketio "github.com/smart-kf/go-socket.io"
	"github.com/smart-kf/go-socket.io/engineio"
	"github.com/smart-kf/go-socket.io/engineio/transport"
	"github.com/smart-kf/go-socket.io/engineio/transport/polling"
	"github.com/smart-kf/go-socket.io/engineio/transport/websocket"

	websocket2 "goim3/application/websocket"
	"goim3/config"
	"goim3/pkg/utils"
)

type WebsocketServer struct {
	socketServer *socketio.Server

	mu              sync.Mutex // protect follow
	conns           map[string]socketio.Conn
	tokenSessionMap map[string]string // token: sessionId
}

func CreateWebsocketServer() *WebsocketServer {
	wsServer := &WebsocketServer{
		conns:           make(map[string]socketio.Conn),
		tokenSessionMap: make(map[string]string),
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
				polling.Default,
				&websocket.Transport{
					ReadBufferSize:   config.Config.SocketIO.ReadBufferSize,
					WriteBufferSize:  config.Config.SocketIO.WriteBufferSize,
					HandshakeTimeout: config.Config.SocketIO.ReadTimeout.Duration(),
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
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
	// auth check.
	if !config.Config.AuthCheck.Enable {
		return nil, nil
	}

	if err := s.doAuthCheck(r); err != nil {
		return nil, err
	}

	return http.Header{}, nil
}

type WebsocketAuthRequest struct {
	Token     string `json:"token"`
	Ip        string `json:"ip"`
	Platform  string `json:"platform"`
	UserAgent string `json:"userAgent"`
}

func (s *WebsocketServer) doAuthCheck(r *http.Request) error {
	client := resty.NewWithClient(
		&http.Client{
			Timeout: config.Config.AuthCheck.Timeout.Duration(),
		},
	)
	authRequest := WebsocketAuthRequest{
		Token:     r.URL.Query().Get("token"),
		Ip:        r.RemoteAddr,
		Platform:  r.URL.Query().Get("platform"),
		UserAgent: r.UserAgent(),
	}
	req := client.R().SetBody(
		authRequest,
	)
	rsp, err := req.Post(config.Config.AuthCheck.HttpUrl)
	if err != nil {
		return err
	}
	rsp.RawBody().Close()
	if rsp.StatusCode() != config.Config.AuthCheck.ResponseCode {
		return errors.New("auth failed")
	}
	return nil
}

// onConnectionInit 验证通过以后得初始化操作.
func (s *WebsocketServer) onConnectionInit(r *http.Request, conn engineio.Conn) {
	return
}

func (s *WebsocketServer) RegisterEvents() {
	s.socketServer.OnConnect(
		"/", s.OnConnect,
	)
	s.socketServer.OnEvent("/", "message", s.onMessage)
	s.socketServer.OnDisconnect("/", s.onDisconnect)
	s.socketServer.OnError(
		"/", func(conn socketio.Conn, err error) {
			s.onDisconnect(conn, err.Error())
		},
	)
}

func (s *WebsocketServer) onMessage(conn socketio.Conn, msg string) {
	var app websocket2.WebsocketApplication
	app.OnMessage(conn.Context().(context.Context), msg)
}

func (s *WebsocketServer) OnConnect(conn socketio.Conn) error {
	u := conn.URL()
	token := u.Query().Get("token")
	platform := u.Query().Get("platform")

	ctx := utils.SetConnContext(context.Background(), token, conn.ID(), platform)
	conn.SetContext(ctx)

	s.mu.Lock()
	s.conns[conn.ID()] = conn
	s.tokenSessionMap[token] = conn.ID()
	s.mu.Unlock()
	var app websocket2.ConnectionApplication
	if err := app.OnConnect(context.Background(), token, conn.ID(), platform); err != nil {
		xlogger.Error(context.Background(), "OnConnect-failed", xlogger.Err(err))
		return err
	} else {
		xlogger.Info(context.Background(), "OnConnect-success", xlogger.Any(token, conn.ID()))
	}
	return nil
}

func (s *WebsocketServer) onDisconnect(conn socketio.Conn, msg string) {
	ctx := utils.MustGetConnContext(conn.Context().(context.Context))
	s.mu.Lock()
	delete(s.conns, ctx.SessionId)
	delete(s.tokenSessionMap, ctx.SessionId)
	s.mu.Unlock()

	var app websocket2.ConnectionApplication
	if err := app.OnDisConnect(context.Background(), ctx.Token, ctx.SessionId, ctx.Platform); err != nil {
		xlogger.Error(context.Background(), "OnDisConnect-failed", xlogger.Err(err))
	} else {
		xlogger.Info(context.Background(), "OnDisConnect-success", xlogger.Any(ctx.Token, conn.ID()))
	}
}

type PushMessageRequest struct {
	SessionIds []string `json:"sessionIds"`
	SessionId  string   `json:"sessionId"`
	Event      string   `json:"event"`
	Data       string   `json:"data"`
}

func (s *WebsocketServer) Push(ctx *gin.Context) {
	var req PushMessageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(200, gin.H{"err": err})
		return
	}
	var sessionIds = req.SessionIds
	if req.SessionId != "" {
		sessionIds = append(sessionIds, req.SessionId)
	}
	var conns []socketio.Conn
	s.mu.Lock()
	for _, sessionId := range sessionIds {
		conn, ok := s.conns[sessionId]
		if ok {
			conns = append(conns, conn)
		}
	}
	s.mu.Unlock()
	if len(conns) == 0 {
		ctx.JSON(
			200, gin.H{
				"success": true,
			},
		)
		return
	}
	for _, conn := range conns {
		conn.Emit(req.Event, req.Data)
	}
	ctx.JSON(
		200, gin.H{
			"success": true,
		},
	)
}
