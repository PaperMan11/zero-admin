// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"zero-admin/api/admin/internal/config"
	"zero-admin/api/admin/internal/middleware"
)

type ServiceContext struct {
	Config           config.Config
	VerifyPermission rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		VerifyPermission: middleware.NewVerifyPermissionMiddleware().Handle,
	}
}
