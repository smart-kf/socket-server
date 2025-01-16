package endpoints

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"goim3/config"
	"goim3/endpoints/network"
)

type HttpServer struct {
	WebsocketServer *network.WebsocketServer
}

func (s *HttpServer) Start() {
	s.WebsocketServer = network.CreateWebsocketServer()

	g := gin.Default()
	g.Use(cors.Default())
	g.Any("/socket.io/", gin.WrapH(s.WebsocketServer))
	g.Static("/static", "./asset")

	// TODO:: push - token
	g.POST("/api/push", s.WebsocketServer.Push)

	wsGroup := g.Group("/api/ws")
	{
		wsGroup.GET("/connections", s.WebsocketServer.Connections)
	}

	http.ListenAndServe(config.Config.ListenAddress, g)
}
