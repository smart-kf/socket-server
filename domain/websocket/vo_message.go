package websocket

type MessageVo struct {
	MsgType     string // text || image || video
	MsgId       string
	GuestName   string // 客户名称
	GuestAvatar string // 客户头像
	KfName      string // 客服名称
	KfAvatar    string // 客服头像
	Content     string // 具体消息内容
	Ip          string // 客户IP
	IsFromKf    bool   // 是否是来自客服发送
}

type VoSessionId struct {
	SessionId string
}
