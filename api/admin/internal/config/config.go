// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"zero-admin/pkg/orm"
)

type Config struct {
	rest.RestConf
	Redis redis.RedisConf
	Mysql orm.Config
	Auth  struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
		ExcludeUrl    []string `json:",optional"`
	}

	// 系统
	SysRpc zrpc.RpcClientConf
}
