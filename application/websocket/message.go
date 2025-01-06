package websocket

import (
	"context"
	"encoding/json"

	"goim3/application/converter"
	"goim3/domain/service/websocket"
	"goim3/pkg/utils"
)

type WebsocketApplication struct{}

func (w *WebsocketApplication) OnMessage(ctx context.Context, messageBody string) {
	connCtx := utils.MustGetConnContext(ctx)
	var dto converter.MessageDTO
	err := json.Unmarshal([]byte(messageBody), &dto)
	if err != nil {
		return
	}
	err = websocket.CreateMessage(ctx, dto.ToModel(connCtx))
	if err != nil {
		return
	}
}
