package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"goim3/config"
)

var Client *redis.Client

func InitRedisClient() {
	client := redis.NewClient(
		&redis.Options{
			Addr:       config.Config.Redis.Address,
			ClientName: "socket.io",
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				if err := cn.Ping(ctx).Err(); err != nil {
					return err
				}
				return nil
			},
			Password:     config.Config.Redis.Password,
			DB:           config.Config.Redis.DB,
			DialTimeout:  config.Config.Redis.Timeout.Duration(),
			ReadTimeout:  config.Config.Redis.Timeout.Duration(),
			WriteTimeout: config.Config.Redis.Timeout.Duration(),
			PoolSize:     config.Config.Redis.PoolSize,
			MinIdleConns: config.Config.Redis.MinIdleConn,
			MaxIdleConns: config.Config.Redis.MaxIdleConn,
		},
	)

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	Client = client
}
