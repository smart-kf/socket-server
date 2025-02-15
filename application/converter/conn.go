package converter

import "goim3/domain/websocket/model"

type ConnDTO struct {
	Token     string `json:"token"`
	SessionId string `json:"sessionId"`
	Platform  string `json:"platform"`
}

func (c *ConnDTO) ToModel() *model.Conn {
	return &model.Conn{
		Token:     c.Token,
		SessionId: c.SessionId,
		Platform:  c.Platform,
	}
}
