package websocket

import (
	"context"

	"goim3/domain/websocket"
	"goim3/domain/websocket/model"
)

func CreateConn(ctx context.Context, conn *model.Conn) error {
	connAgg := websocket.FactoryConnAgg(ctx, conn)
	return connAgg.Create(ctx)
}

func DeleteConn(ctx context.Context, conn *model.Conn) error {
	connAgg := websocket.FactoryConnAgg(ctx, conn)
	return connAgg.Delete(ctx)
}
