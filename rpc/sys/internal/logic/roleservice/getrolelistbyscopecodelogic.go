package roleservicelogic

import (
	"context"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleListByScopeCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleListByScopeCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListByScopeCodeLogic {
	return &GetRoleListByScopeCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取关联给定安全范围的角色
func (l *GetRoleListByScopeCodeLogic) GetRoleListByScopeCode(in *sysclient.GetRolesByScopeCodeRequest) (*sysclient.GetRoleByRoleCodesResponse, error) {
	res, _ := l.svcCtx.DB.GetRolePermsByScopeCode(l.ctx, in.ScopeCode)
	return &sysclient.GetRoleByRoleCodesResponse{
		RolePerms: logic.ConvertToRpcRoleScopes(res),
	}, nil
}
