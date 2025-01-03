package websocket

import (
	"context"
	"errors"

	"github.com/alibaba/ioc-golang/autowire/singleton"

	"goim3/domain/websocket/gateway"
	"goim3/domain/websocket/model"
)

type MessageAgg struct {
	message        *model.Message
	messageGateway gateway.MessageGateway
}

func (a *MessageAgg) SetMessage(message *model.Message) {
	a.message = message
}

func (a *MessageAgg) Create(ctx context.Context) error {
	if err := a.init(); err != nil {
		return err
	}
	return a.messageGateway.Create(ctx, a.message)
}

func (a *MessageAgg) init() error {
	if a.messageGateway == nil {
		msgGateway, err := singleton.GetImpl("MessageGateway", nil)
		if err != nil {
			panic(err)
		}
		a.messageGateway = msgGateway.(gateway.MessageGateway)
	}

	if a.message == nil {
		return errors.New("message can't be nil")
	}

	return nil
}
