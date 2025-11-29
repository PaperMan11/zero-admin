package interceptor

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"zero-admin/rpc/sys/db"
)

const (
	OperatorIDKey = "operator_id"
)

// 用户权限拦截器
type PermissionInterceptor struct {
	db db.DB
}

func NewPermissionInterceptor(db db.DB) *PermissionInterceptor {
	return &PermissionInterceptor{
		db: db,
	}
}

func (p *PermissionInterceptor) VerifyPermission(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logc.Errorf(ctx, "no metadata")
		return nil, status.Error(codes.InvalidArgument, "no metadata")
	}
	logc.Debugf(ctx, "incoming metadata: %v", incomingContext)
	logc.Debugf(ctx, "info: %+v", info)
	return handler(ctx, req)
}
