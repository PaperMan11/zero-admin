// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"zero-admin/api/admin/internal/config"
	"zero-admin/api/admin/internal/middleware"
	casbinUtil "zero-admin/pkg/casbin"
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
	CasbinEnforcer   *casbin.SyncedCachedEnforcer

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

	// casbin
	enforcer := casbinUtil.MustNewCasbinEnforcer("", casbinUtil.MustNewGormAdapter(c.Mysql.DSN))
	wather := casbinUtil.MustNewRedisWatcher(&c.Redis, func(data string) {
		logx.Infof("casbin watcher: %s", data)
	})
	enforcer.SetWatcher(wather)
	enforcer.EnableAutoSave(true)

	return &ServiceContext{
		Config:         c,
		Redis:          redis.MustNewRedis(c.Redis),
		CasbinEnforcer: enforcer,

		VerifyPermission: middleware.NewVerifyPermissionMiddleware(enforcer, c.Auth.ExcludeUrl...).Handle,
		AddLog:           middleware.NewAddLogMiddleware(operateLogService).Handle,

		// 系统相关
		AuthService:       authservice.NewAuthService(sysClient),
		OperateLogService: operateLogService,
		ScopeService:      scopeservice.NewScopeService(sysClient),
		RoleService:       roleservice.NewRoleService(sysClient),
		UserService:       userservice.NewUserService(sysClient),
	}
}
