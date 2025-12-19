// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/scopeservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnassignedMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUnassignedMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnassignedMenusLogic {
	return &GetUnassignedMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnassignedMenusLogic) GetUnassignedMenus(req *types.Empty) (resp *types.GetUnassignedMenusResponse, err error) {
	res, _ := l.svcCtx.ScopeService.GetUnassignedMenus(l.ctx, &scopeservice.Empty{})
	return &types.GetUnassignedMenusResponse{
		Menus: utils.ConvertToTypesMenus(res.Menus),
	}, nil
}
