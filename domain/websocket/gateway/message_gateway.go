package gateway

import (
	"context"

	"goim3/domain/websocket/model"
)

type MessageGateway interface {
	Create(ctx context.Context, message *model.Message) error
}
