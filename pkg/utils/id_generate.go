package utils

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"os"
	"sync"
	"time"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	"github.com/google/uuid"
	"github.com/sony/sonyflake"
)

var MachineID uint16 = 1

func init() {
	singleton.RegisterStructDescriptor(
		&autowire.StructDescriptor{
			Factory: func() interface{} {
				return NewUUIDGenerate()
			},
			Alias: "IDGenerator",
		},
	)
}

var (
	hostname string
)

func init() {
	hostname, _ = os.Hostname()
}

func NewUUIDGenerate() *IDGenerator {
	return &IDGenerator{
		hash: sha256.New(),
		sonyFlake: sonyflake.NewSonyflake(
			sonyflake.Settings{
				StartTime: time.Now(),
				MachineID: func() (uint16, error) {
					return MachineID, nil
				},
				CheckMachineID: nil,
			},
		),
	}
}

// IDGenerator 2种id生成方式
// string: 基于 hostname + uuid
// int64: 基于 snowflake
type IDGenerator struct {
	mu   sync.Mutex
	hash hash.Hash

	sonyFlake *sonyflake.Sonyflake
}

func (u *IDGenerator) NewID() string {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.hash.Write([]byte(hostname + uuid.New().String()))
	hashString := u.hash.Sum(nil)
	u.hash.Reset()
	return fmt.Sprintf("%x", hashString)
}

func (u *IDGenerator) NewIdInt64String() int64 {
	id, _ := u.sonyFlake.NextID()
	return int64(id)
}
