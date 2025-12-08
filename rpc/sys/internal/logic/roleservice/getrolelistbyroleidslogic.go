package roleservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleListByRoleIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleListByRoleIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListByRoleIDsLogic {
	return &GetRoleListByRoleIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取角色列表
func (l *GetRoleListByRoleIDsLogic) GetRoleListByRoleIDs(in *sysclient.GetRoleByRoleCodesRequest) (*sysclient.GetRoleByRoleCodesResponse, error) {
	res, err := l.svcCtx.DB.GetRoleByCodes(l.ctx, in.RoleCodes)
	if err != nil {
		logc.Errorf(l.ctx, "查询角色列表异常：%v", err)
		return nil, status.Error(codes.Internal, "查询角色列表异常")
	}

	return &sysclient.GetRoleByRoleCodesResponse{
		Roles: logic.ConvertToRpcRoles(res),
	}, nil
}
