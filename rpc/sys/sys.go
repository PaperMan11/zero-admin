package main

import (
	"flag"
	"fmt"

	"zero-admin/rpc/sys/internal/config"
	authserviceServer "zero-admin/rpc/sys/internal/server/authservice"
	roleserviceServer "zero-admin/rpc/sys/internal/server/roleservice"
	scopeserviceServer "zero-admin/rpc/sys/internal/server/scopeservice"
	userserviceServer "zero-admin/rpc/sys/internal/server/userservice"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/sys.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		sysclient.RegisterAuthServiceServer(grpcServer, authserviceServer.NewAuthServiceServer(ctx))
		sysclient.RegisterRoleServiceServer(grpcServer, roleserviceServer.NewRoleServiceServer(ctx))
		sysclient.RegisterScopeServiceServer(grpcServer, scopeserviceServer.NewScopeServiceServer(ctx))
		sysclient.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(ctx.PermissionInterceptor.VerifyPermission)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
