package main

import (
	"flag"

	xlogger "github.com/clearcodecn/log"

	"goim3/config"
	"goim3/endpoints"
	_ "goim3/infrastructure/gateway_impl"
	"goim3/infrastructure/nsq"
	"goim3/infrastructure/redis"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "config.yaml", "配置文件")
}

func main() {
	flag.Parse()

	config.Load(configFile)
	nsq.InitProducer()
	redis.InitRedisClient()

	var HttpServer endpoints.HttpServer

	HttpServer.Start()
}

func initLogger() {
	// xlogger.AddHook(func(ctx context.Context) xlogger.Field {
	//	reqid, ok := ctx.Value("reqid").(string)
	//	if !ok {
	//		return xlogger.Field{}
	//	}
	//	return xlogger.Any("reqid", reqid)
	// })
	logger, err := xlogger.NewLog(
		xlogger.Config{
			Level:  config.Config.Log.Level,
			Format: config.Config.Log.Format,
			File:   config.Config.Log.File,
		},
	)
	if err != nil {
		panic(err)
	}

	xlogger.SetGlobal(logger)
}
