package dao

type Conn struct {
	SessionId string `json:"sessionId"`
	Token     string `json:"token"`
	Platform  string `json:"platform"`
}
