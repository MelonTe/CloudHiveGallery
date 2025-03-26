package redlock

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

//引入redis分布式锁实现redsync，解决误删（值校验）、分布式锁丢失问题（redis distribution lock）

var rs *redsync.Redsync

// 初始化分布式锁，被redis包的init函数调用
func InitRedSync(redClient *redis.Client) {
	rs = redsync.New(goredis.NewPool(redClient))
}

func GetRedSync() *redsync.Redsync {
	return rs
}
