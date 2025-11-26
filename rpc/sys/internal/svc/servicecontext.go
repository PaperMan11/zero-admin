package svc

import (
	"zero-admin/rpc/sys/db"
	"zero-admin/rpc/sys/interceptor"
	"zero-admin/rpc/sys/internal/config"
)

type ServiceContext struct {
	Config                config.Config
	DB                    db.DB
	PermissionInterceptor *interceptor.PermissionInterceptor
}

func NewServiceContext(c config.Config) *ServiceContext {
	_db := db.MustNewDB(c.DbMode, &c.Mysql)
	return &ServiceContext{
		Config:                c,
		DB:                    _db,
		PermissionInterceptor: interceptor.NewPermissionInterceptor(_db),
	}
}
