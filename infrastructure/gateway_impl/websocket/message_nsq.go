package websocket

import (
	"context"
	"encoding/json"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	xlogger "github.com/clearcodecn/log"

	"goim3/config"
	"goim3/domain/websocket/model"
	"goim3/infrastructure/nsq"
	"goim3/infrastructure/nsq/dao/websocket"
)

func init() {
	singleton.RegisterStructDescriptor(
		&autowire.StructDescriptor{
			Factory: func() interface{} {
				return &MessageNsqImpl{}
			},
			Alias: "MessageGateway",
		},
	)
}

type MessageNsqImpl struct {
}

func (m *MessageNsqImpl) Create(ctx context.Context, message *model.Message) error {
	messageDao := websocket.Message{
		Event:       message.Event,
		MsgType:     message.MsgType,
		MsgId:       message.MsgId,
		GuestName:   message.GuestName,
		GuestAvatar: message.GuestAvatar,
		KfName:      message.KfName,
		KfAvatar:    message.KfAvatar,
		Content:     message.Content,
		Ip:          message.Ip,
		Platform:    message.Platform,
		SessionId:   message.SessionId,
		Token:       message.Token,
	}
	body, _ := json.Marshal(messageDao)

	xlogger.Info(ctx, "nsq produce message", xlogger.Any("msg", string(body)))
	if err := nsq.NSQProducer.Publish(config.Config.Nsq.MessageTopic, body); err != nil {
		return err
	}
	return nil
}
