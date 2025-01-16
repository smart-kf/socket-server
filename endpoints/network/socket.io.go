package network

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

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

	connTotal       int64      // 启动服务器至今的连接数量
	authFailedTotal int64      // auth 失败的连接
	msgTotal        int64      // 推送消息的数量
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

	go func() {
		tk := time.NewTicker(60 * time.Second)
		for range tk.C {
			wsServer.getWsResult()
		}
	}()

	return wsServer
}

func (s *WebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.socketServer.ServeHTTP(w, r)
}

// onRequestCheck 建立连接的验证代码
func (s *WebsocketServer) onRequestCheck(r *http.Request) (http.Header, error) {
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
		atomic.AddInt64(&s.authFailedTotal, 1)
		return err
	}
	rsp.RawBody().Close()
	if rsp.StatusCode() != config.Config.AuthCheck.ResponseCode {
		atomic.AddInt64(&s.authFailedTotal, 1)
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

	atomic.AddInt64(&s.connTotal, 1)

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

		fmt.Println("send message -->", string(req.Data))
	}

	atomic.AddInt64(&s.msgTotal, int64(len(conns)))
	ctx.JSON(
		200, gin.H{
			"success": true,
		},
	)
}

type WsConnResult struct {
	ConnTotal          int64 `json:"conn_total"`
	MsgTotal           int64 `json:"msg_total"`
	CurrentConnections int64 `json:"current_connections"`
	AuthFailedTotal    int64 `json:"auth_failed_total"`
}

func (s *WebsocketServer) Connections(ctx *gin.Context) {
	ctx.JSON(200, s.getWsResult())
}

func (s *WebsocketServer) getWsResult() WsConnResult {
	var rsp WsConnResult
	rsp.MsgTotal = atomic.LoadInt64(&s.msgTotal)
	rsp.ConnTotal = atomic.LoadInt64(&s.connTotal)
	rsp.AuthFailedTotal = atomic.LoadInt64(&s.authFailedTotal)
	s.mu.Lock()
	rsp.CurrentConnections = int64(len(s.conns))
	s.mu.Unlock()

	log.Printf(
		"[连接信息]: 服务启动至今总连接=%d  当前活跃连接:%d 推送消息数量: %d auth失败: %d \n",
		rsp.ConnTotal,
		rsp.CurrentConnections,
		rsp.MsgTotal,
		rsp.AuthFailedTotal,
	)
	return rsp
}
