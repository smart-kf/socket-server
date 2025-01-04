package websocket

import (
	"context"

	"goim3/application/converter"
	"goim3/domain/service/websocket"
)

type ConnectionApplication struct{}

// OnConnect 创建一个websocket连接.
func (a *ConnectionApplication) OnConnect(ctx context.Context, token string, sessionId string) error {
	dto := converter.ConnDTO{
		Token:     token,
		SessionId: sessionId,
	}

	return websocket.CreateConn(ctx, dto.ToModel())
}

// OnDisConnect 删除一个websocket连接
func (a *ConnectionApplication) OnDisConnect(ctx context.Context, token string, sessionId string) error {
	dto := converter.ConnDTO{
		Token:     token,
		SessionId: sessionId,
	}

	return websocket.DeleteConn(ctx, dto.ToModel())
}
