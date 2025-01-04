package gateway

import (
	"context"

	"goim3/domain/websocket/model"
)

type ConnGateway interface {
	Create(ctx context.Context, conn *model.Conn) error
	Delete(ctx context.Context, conn *model.Conn) error
}
