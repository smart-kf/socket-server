package websocket

import (
	"context"
	"errors"

	"github.com/alibaba/ioc-golang/autowire/singleton"

	"goim3/domain/websocket/gateway"
	"goim3/domain/websocket/model"
)

type ConnAgg struct {
	conn        *model.Conn
	connGateway gateway.ConnGateway
}

func (a *ConnAgg) SetConn(conn *model.Conn) {
	a.conn = conn
}

func (a *ConnAgg) init() error {
	if a.conn == nil {
		return errors.New("connection is not init")
	}
	if a.conn.SessionId == "" || a.conn.Token == "" {
		return errors.New("connection sessionId or token is empty")
	}
	if a.connGateway == nil {
		connGateway, err := singleton.GetImpl("ConnGateway", nil)
		if err != nil {
			panic(err)
		}
		a.connGateway = connGateway.(gateway.ConnGateway)
	}
	return nil
}

func (a *ConnAgg) Delete(ctx context.Context) error {
	if err := a.init(); err != nil {
		return err
	}
	return a.connGateway.Delete(ctx, a.conn)
}

func (a *ConnAgg) Create(ctx context.Context) error {
	if err := a.init(); err != nil {
		return err
	}
	return a.connGateway.Create(ctx, a.conn)
}
