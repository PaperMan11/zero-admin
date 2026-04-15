// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"zero-admin/pkg/orm"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
	}
	Redis redis.RedisConf
	Mysql orm.Config
}
