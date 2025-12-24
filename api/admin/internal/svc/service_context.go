// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
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
	LocalCache       *collection.Cache
	Redis            *redis.Redis
	VerifyPermission rest.Middleware
	AddLog           rest.Middleware
	JwtExpireAuth    rest.Middleware

	Barrier syncx.SingleFlight
	// casbin
	CasbinEnforcer *casbin.SyncedCachedEnforcer

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

	// singleFlight
	barrier := syncx.NewSingleFlight()

	// cache
	redisCli := redis.MustNewRedis(c.Redis)
	localCache, err := collection.NewCache(time.Minute*30, collection.WithName("cache"))
	if err != nil {
		logx.Must(err)
	}

	// casbin
	enforcer := casbinUtil.MustNewCasbinEnforcer("", casbinUtil.MustNewGormAdapter(c.Mysql.DSN))
	wather := casbinUtil.MustNewRedisWatcher(&c.Redis, func(data string) {
		logx.Debugf("casbin watcher: %s", data)
		enforcer.LoadPolicy()
	})
	err = enforcer.SetWatcher(wather)
	if err != nil {
		logx.Errorf("set casbin watcher error: %v", err)
	}
	enforcer.EnableAutoSave(true)

	return &ServiceContext{
		Config:         c,
		Redis:          redisCli,
		LocalCache:     localCache,
		CasbinEnforcer: enforcer,
		Barrier:        barrier,

		JwtExpireAuth:    middleware.NewJwtExpireAuthMiddleware(redisCli, localCache).Handle,
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
