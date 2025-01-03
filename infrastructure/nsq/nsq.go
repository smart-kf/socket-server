package nsq

import (
	"os"
	"time"

	"github.com/nsqio/go-nsq"

	"goim3/config"
)

var NSQProducer *nsq.Producer

func InitProducer() {
	hostname, _ := os.Hostname()
	cfg := nsq.NewConfig()
	cfg.DialTimeout = 60 * time.Second
	cfg.ReadTimeout = 60 * time.Second
	cfg.WriteTimeout = 60 * time.Second
	cfg.ClientID = hostname
	cfg.Hostname = hostname + "-websocketserver"
	cfg.UserAgent = "go-" + hostname + "-websocketserver"
	p, err := nsq.NewProducer(config.Config.Nsq.Address, cfg)
	if err != nil {
		panic(err)
	}
	NSQProducer = p
}
