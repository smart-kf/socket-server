package model

// Message 发送的消息.
type Message struct {
	Event       string
	MsgType     string // text || image || video
	MsgId       string
	GuestName   string // 客户名称
	GuestAvatar string // 客户头像
	KfName      string // 客服名称
	KfAvatar    string // 客服头像
	Content     string // 具体消息内容
	Ip          string // 客户IP
	Platform    string `json:"platform"`
	SessionId   string `json:"sessionId"` // sessionId
	Token       string `json:"token"`
}
