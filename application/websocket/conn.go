package websocket

import (
	"context"

	"goim3/application/converter"
	"goim3/domain/common/constant"
	"goim3/domain/service/websocket"
	"goim3/domain/websocket/model"
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

	if err := websocket.CreateMessage(
		ctx, &model.Message{
			SessionId: sessionId,
			Platform:  platform,
			Event:     constant.EventSessionId,
			Token:     token,
		},
	); err != nil {
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

	if err := websocket.CreateMessage(
		ctx, &model.Message{
			SessionId: sessionId,
			Platform:  platform,
			Event:     constant.EventSessionId,
			Token:     token,
		},
	); err != nil {
		return err
	}

	return websocket.DeleteConn(ctx, dto.ToModel())
}
