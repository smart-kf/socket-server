package websocket

import (
	"context"

	"goim3/domain/websocket/model"
)

func FactoryMessageAgg(ctx context.Context, msg *model.Message) *MessageAgg {
	var agg = &MessageAgg{}
	agg.SetMessage(msg)

	return agg
}

func FactoryConnAgg(ctx context.Context, conn *model.Conn) *ConnAgg {
	var agg = &ConnAgg{}
	agg.SetConn(conn)

	return agg
}
