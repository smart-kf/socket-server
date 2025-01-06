package websocket

import (
	"context"

	"goim3/domain/websocket"
	"goim3/domain/websocket/model"
)

func CreateMessage(ctx context.Context, message *model.Message) error {
	agg := websocket.FactoryMessageAgg(ctx, message)

	return agg.Create(ctx)
}
