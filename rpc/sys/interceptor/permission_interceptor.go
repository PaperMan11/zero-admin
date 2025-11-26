package interceptor

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"zero-admin/pkg/convert"
	"zero-admin/rpc/sys/db"
	perm "zero-admin/rpc/sys/db/common"
	"zero-admin/rpc/sys/db/mysql/model"
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
	operator := incomingContext.Get(OperatorIDKey)
	if len(operator) == 0 {
		logc.Errorf(ctx, "no operator")
		return nil, status.Error(codes.InvalidArgument, "no operator")
	}
	operatorID := convert.ToInt64(operator[0])
	roleCodes, _ := p.db.GetUserRoleCodes(ctx, operatorID)
	scopePerms, _ := p.db.GetRolesScopesPerm(ctx, roleCodes)

	logc.Debugf(ctx, "info: %+v", info)
	if !hasScopePermission(scopePerms, "", perm.PERM_READ) {
		logc.Errorf(ctx, "no scope permission: operator: %d", operatorID)
		return nil, status.Error(codes.PermissionDenied, "no scope permission")
	}
	return handler(ctx, req)
}

// 是否有scope的权限
func hasScopePermission(scopePerms []model.SysRoleScope, scopeCode string, perm perm.PermType) bool {
	for _, sp := range scopePerms {
		if sp.ScopeCode == scopeCode && sp.Perm&perm > 0 {
			return true
		}
	}
	return false
}
