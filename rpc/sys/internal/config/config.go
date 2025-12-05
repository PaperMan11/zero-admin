package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zero-admin/pkg/orm"
)

type Config struct {
	zrpc.RpcServerConf

	DbMode string `json:",default=mockdb"` // 数据库模式 mysql, mockdb
	Mysql  orm.Config

	Jwt struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
	}
}
