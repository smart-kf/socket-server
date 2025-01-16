package websocket

import (
	"context"

	"goim3/application/converter"
	"goim3/domain/service/websocket"
)

type ConnectionApplication struct{}

// OnConnect 创建一个websocket连接.
func (a *ConnectionApplication) OnConnect(
	ctx context.Context,
	token string,
	sessionId string,
	platform string,
) error {
	dto := converter.ConnDTO{
		Token:     token,
		SessionId: sessionId,
		Platform:  platform,
	}

	if err := websocket.CreateConn(ctx, dto.ToModel()); err != nil {
		return err
	}

	return nil
}

// OnDisConnect 删除一个websocket连接
func (a *ConnectionApplication) OnDisConnect(
	ctx context.Context,
	token string,
	sessionId string,
	platform string,
) error {
	dto := converter.ConnDTO{
		Token:     token,
		SessionId: sessionId,
		Platform:  platform,
	}

	return websocket.DeleteConn(ctx, dto.ToModel())
}
