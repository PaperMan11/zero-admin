package roleservicelogic

import (
	"context"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllRolesLogic {
	return &GetAllRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 全部角色
func (l *GetAllRolesLogic) GetAllRoles(in *sysclient.GetAllRolesRequest) (*sysclient.GetAllRolesResponse, error) {
	roles, _ := l.svcCtx.DB.GetAllRoles(l.ctx)
	return &sysclient.GetAllRolesResponse{
		Roles: logic.ConvertToRpcRoles(roles),
	}, nil
}
