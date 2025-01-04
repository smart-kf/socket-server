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

var (
	configFile string
)

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

	// server := socketio.NewServer(nil)
	//
	// server.OnConnect(
	// 	"/", func(s socketio.Conn) error {
	// 		s.SetContext("")
	// 		fmt.Println("connected:", s.ID())
	// 		return nil
	// 	},
	// )
	//
	// server.OnEvent(
	// 	"/", "notice", func(s socketio.Conn, msg string) {
	// 		fmt.Println("notice:", msg)
	// 		s.Emit("reply", "have "+msg)
	// 	},
	// )
	//
	// server.OnEvent(
	// 	"/chat", "msg", func(s socketio.Conn, msg string) string {
	// 		s.SetContext(msg)
	// 		return "recv " + msg
	// 	},
	// )
	//
	// server.OnEvent(
	// 	"/", "bye", func(s socketio.Conn) string {
	// 		last := s.Context().(string)
	// 		s.Emit("bye", last)
	// 		s.Close()
	// 		return last
	// 	},
	// )
	//
	// server.OnError(
	// 	"/", func(s socketio.Conn, e error) {
	// 		fmt.Println("meet error:", e)
	// 	},
	// )
	//
	// server.OnDisconnect(
	// 	"/", func(s socketio.Conn, reason string) {
	// 		fmt.Println("closed", reason)
	// 	},
	// )
	//
	// go server.Serve()
	// defer server.Close()
	//
	// http.Handle("/socket.io/", server)
	// http.Handle("/", http.FileServer(http.Dir("./asset")))
	// log.Println("Serving at localhost:8000...")
	// log.Fatal(http.ListenAndServe(":8000", nil))
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
