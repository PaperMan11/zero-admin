// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"zero-admin/api/admin/internal/config"
	"zero-admin/api/admin/internal/middleware"
	"zero-admin/rpc/sys/client/authservice"
	"zero-admin/rpc/sys/client/operatelogservice"
	"zero-admin/rpc/sys/client/roleservice"
	"zero-admin/rpc/sys/client/scopeservice"
	"zero-admin/rpc/sys/client/userservice"
)

type ServiceContext struct {
	Config           config.Config
	Redis            *redis.Redis
	VerifyPermission rest.Middleware
	AddLog           rest.Middleware

	// 系统相关
	AuthService       authservice.AuthService
	OperateLogService operatelogservice.OperateLogService
	ScopeService      scopeservice.ScopeService
	RoleService       roleservice.RoleService
	UserService       userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	sysClient := zrpc.MustNewClient(c.SysRpc)
	operateLogService := operatelogservice.NewOperateLogService(sysClient)
	return &ServiceContext{
		Config:           c,
		Redis:            redis.MustNewRedis(c.Redis),
		VerifyPermission: middleware.NewVerifyPermissionMiddleware().Handle,
		AddLog:           middleware.NewAddLogMiddleware(operateLogService).Handle,

		// 系统相关
		AuthService:       authservice.NewAuthService(sysClient),
		OperateLogService: operateLogService,
		ScopeService:      scopeservice.NewScopeService(sysClient),
		RoleService:       roleservice.NewRoleService(sysClient),
		UserService:       userservice.NewUserService(sysClient),
	}
}
