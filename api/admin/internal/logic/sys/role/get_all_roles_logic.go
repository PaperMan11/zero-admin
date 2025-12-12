// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/sysclient"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllRolesLogic {
	return &GetAllRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllRolesLogic) GetAllRoles(req *types.GetAllRolesRequest) (resp *types.GetAllRolesResponse, err error) {
	res, _ := l.svcCtx.RoleService.GetAllRoles(l.ctx, &sysclient.GetAllRolesRequest{})
	return &types.GetAllRolesResponse{
		Roles: utils.ConvertToTypesRoles(res.Roles),
	}, nil
}
