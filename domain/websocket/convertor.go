package websocket

import "goim3/domain/websocket/model"

func Message2Model(vo MessageVo) *model.Message {
	return &model.Message{
		MsgType:     vo.MsgType,
		MsgId:       vo.MsgId,
		GuestName:   vo.GuestName,
		GuestAvatar: vo.GuestAvatar,
		KfName:      vo.KfName,
		KfAvatar:    vo.KfAvatar,
		Content:     vo.Content,
		Ip:          vo.Ip,
		IsFromKf:    vo.IsFromKf,
	}
}
