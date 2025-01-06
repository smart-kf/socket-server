package converter

import "goim3/domain/websocket/model"

type MessageDTO struct {
	MsgType     string `json:"msgType"`     // text || image || video
	MsgId       string `json:"msgId"`       // 消息id
	GuestName   string `json:"guestName"`   // 客户名称
	GuestAvatar string `json:"guestNvatar"` // 客户头像
	KfName      string `json:"kfName"`      // 客服名称
	KfAvatar    string `json:"kfAvatar"`    // 客服头像
	Content     string `json:"content"`     // 具体消息内容
	Ip          string `json:"ip"`          // 客户IP
	Token       string `json:"token"`
	Platform    string `json:"platform"`
	SessionId   string `json:"sessionId"`
}

func (m *MessageDTO) ToModel() *model.Message {
	return &model.Message{
		MsgType:     m.MsgType,
		MsgId:       m.MsgId,
		GuestName:   m.GuestName,
		GuestAvatar: m.GuestAvatar,
		KfName:      m.KfAvatar,
		KfAvatar:    m.KfAvatar,
		Content:     m.Content,
		Ip:          m.Ip,
		Token:       m.Token,
		Platform:    m.Platform,
		SessionId:   m.SessionId,
	}
}
