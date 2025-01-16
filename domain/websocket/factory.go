package websocket

import (
	"context"

	"goim3/domain/websocket/model"
)

func FactoryMessageAgg(ctx context.Context, msg *model.Message) *MessageAgg {
	agg := &MessageAgg{}
	agg.SetMessage(msg)

	return agg
}
