package conn

import (
	"context"
	"fmt"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/singleton"

	"goim3/config"
	"goim3/domain/websocket/model"
	"goim3/infrastructure/redis"
	"goim3/infrastructure/redis/dao"
)

type ConnRedisImpl struct{}

func (ConnRedisImpl) getConnKey() string {
	return fmt.Sprintf("%s.%s", config.Config.Redis.KeyPrefix, "connection.hash")
}

// Create 保存用户连接信息，用户token -> sessionId 的映射, 注意这里存储的不是json数据.
// hash  key= prefix.token  , value = sessionId
func (c *ConnRedisImpl) Create(ctx context.Context, conn *model.Conn) error {
	var connDao = dao.Conn{
		SessionId: conn.SessionId,
		Token:     conn.Token,
	}
	return redis.Client.HSet(ctx, c.getConnKey(), connDao.Token, connDao.SessionId).Err()
}

func (c *ConnRedisImpl) Delete(ctx context.Context, conn *model.Conn) error {
	var connDao = dao.Conn{
		SessionId: conn.SessionId,
		Token:     conn.Token,
	}
	return redis.Client.HDel(ctx, c.getConnKey(), connDao.Token).Err()
}

func init() {
	singleton.RegisterStructDescriptor(
		&autowire.StructDescriptor{
			Factory: func() interface{} {
				return &ConnRedisImpl{}
			},
			Alias: "ConnGateway",
		},
	)
}
