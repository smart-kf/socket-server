package websocket

import (
	"context"

	socketio "github.com/smart-kf/go-socket.io"

	"goim3/application/converter"
	"goim3/domain/service/websocket"
	websocket2 "goim3/domain/websocket"
)

type ConnectionApplication struct{}

// OnConnect 创建一个websocket连接.
func (a *ConnectionApplication) OnConnect(
	ctx context.Context, conn socketio.Conn, token string,
	sessionId string,
) error {
	dto := converter.ConnDTO{
		Token:     token,
		SessionId: sessionId,
	}

	if err := websocket.CreateConn(ctx, dto.ToModel()); err != nil {
		return err
	}

	// 推送连接成功事件
	conn.Emit(
		"connected", websocket2.VoSessionId{
			SessionId: conn.ID(),
		},
	)

	return nil
}

// OnDisConnect 删除一个websocket连接
func (a *ConnectionApplication) OnDisConnect(ctx context.Context, token string, sessionId string) error {
	dto := converter.ConnDTO{
		Token:     token,
		SessionId: sessionId,
	}

	return websocket.DeleteConn(ctx, dto.ToModel())
}
