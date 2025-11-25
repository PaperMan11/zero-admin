package svc

import (
	"zero-admin/rpc/sys/db"
	"zero-admin/rpc/sys/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     db.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     db.MustNewDB(db.DB_MYSQL, &c.Mysql),
	}
}
