package config

import (
	"github.com/make-money-fast/xconfig"

	"goim3/pkg/utils"
)

var Config *config

type config struct {
	MachineID     uint16    `json:"machineId" default:"1"`
	ListenAddress string    `json:"listenAddress" default:":9000"`
	SocketIO      SocketIO  `json:"socketIo"`
	Nsq           Nsq       `json:"nsq"`
	Log           Log       `json:"log"`
	AuthCheck     AuthCheck `json:"authCheck"`
	Redis         Redis     `json:"redis"`
}

type SocketIO struct {
	ReadTimeout     utils.Duration `json:"readTimeout"`
	PingTimeout     utils.Duration `json:"pingTimeout"`
	PingInterval    utils.Duration `json:"pingInterval"`
	ReadBufferSize  int            `json:"readBufferSize"`
	WriteBufferSize int            `json:"writeBufferSize"`
}

type Nsq struct {
	Address      string `json:"address"`
	MessageTopic string `json:"messageTopic"`
}

type AuthCheck struct {
	Enable       bool           `json:"enable"`
	HttpUrl      string         `json:"httpUrl"`      // 检测的url
	ResponseCode int            `json:"responseCode"` // 成功返回的状态码
	Timeout      utils.Duration `json:"timeout"`      // 检测超时.
}

type Redis struct {
	Address     string         `json:"address"`
	DB          int            `json:"db"`
	Password    string         `json:"password"`
	KeyPrefix   string         `json:"keyPrefix"`
	Timeout     utils.Duration `json:"timeout"`
	PoolSize    int            `json:"poolSize"`
	MinIdleConn int            `json:"minIdleConn"`
	MaxIdleConn int            `json:"maxIdleConn"`
}

type Log struct {
	Level  string `json:"level" default:"info"`
	Format string `json:"format" default:"json"`
	File   string `json:"file"`
}

func Load(filename string) {
	var c config
	err := xconfig.ParseFromFile(filename, &c)
	if err != nil {
		panic(err)
	}
	utils.MachineID = c.MachineID
	Config = &c
	return
}
