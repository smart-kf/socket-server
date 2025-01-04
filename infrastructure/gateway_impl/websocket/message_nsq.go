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
	"goim3/pkg/utils"
)

func init() {
	singleton.RegisterStructDescriptor(
		&autowire.StructDescriptor{
			Factory: func() interface{} {
				idgen, err := singleton.GetImpl("IDGenerator", nil)
				if err != nil {
					panic(err)
				}
				return &MessageNsqImpl{
					idGenerator: idgen.(*utils.IDGenerator),
				}
			},
			Alias: "MessageGateway",
		},
	)
}

type MessageNsqImpl struct {
	idGenerator *utils.IDGenerator
}

func (m *MessageNsqImpl) Create(ctx context.Context, message *model.Message) error {
	var messageDao = websocket.Message{
		MsgType:     message.MsgType,
		MsgId:       message.MsgId,
		GuestName:   message.GuestName,
		GuestAvatar: message.GuestAvatar,
		KfName:      message.KfName,
		KfAvatar:    message.KfAvatar,
		Content:     message.Content,
		Ip:          message.Ip,
		IsFromKf:    message.IsFromKf,
	}
	messageDao.MsgId = m.idGenerator.NewID()
	body, _ := json.Marshal(messageDao)

	xlogger.Info(ctx, "nsq produce message", xlogger.Any("msg", string(body)))
	if err := nsq.NSQProducer.Publish(config.Config.Nsq.MessageTopic, body); err != nil {
		return err
	}
	return nil
}
