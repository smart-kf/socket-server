package websocket

import (
	"context"

	"goim3/domain/common/constant"
	"goim3/domain/websocket"
	"goim3/domain/websocket/model"
)

func CreateConn(ctx context.Context, conn *model.Conn) error {
	connAgg := websocket.FactoryConnAgg(ctx, conn)
	if err := connAgg.Create(ctx); err != nil {
		return err
	}
	// 发布一个上线事件.
	agg := websocket.FactoryMessageAgg(
		ctx, &model.Message{
			Event:     constant.Online,
			Platform:  conn.Platform,
			SessionId: conn.SessionId,
			Token:     conn.Token,
		},
	)
	return agg.Create(ctx)
}

func DeleteConn(ctx context.Context, conn *model.Conn) error {
	connAgg := websocket.FactoryConnAgg(ctx, conn)
	if err := connAgg.Delete(ctx); err != nil {
		return err
	}

	// 发布一个离线事件.
	agg := websocket.FactoryMessageAgg(
		ctx, &model.Message{
			Event:     constant.Offline,
			Platform:  conn.Platform,
			SessionId: conn.SessionId,
			Token:     conn.Token,
		},
	)
	return agg.Create(ctx)
}
