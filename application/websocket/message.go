package websocket

import (
	"context"
	"encoding/json"

	"goim3/application/converter"
	"goim3/domain/service/websocket"
)

type WebsocketApplication struct{}

func (w *WebsocketApplication) OnMessage(ctx context.Context, messageBody string) {
	var dto converter.MessageDTO
	err := json.Unmarshal([]byte(messageBody), &dto)
	if err != nil {
		return
	}
	err = websocket.CreateMessage(ctx, dto.ToModel())
	if err != nil {
		return
	}
}
