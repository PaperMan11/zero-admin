// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/logic"
	"zero-admin/rpc/sys/client/scopeservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuTreeLogic) GetMenuTree(req *types.MenuListRequest) (resp *types.MenuTreeResponse, err error) {
	res, err := l.svcCtx.ScopeService.GetMenuTree(l.ctx, &scopeservice.MenuListRequest{
		WithButton: req.WithButton,
		Status:     req.Status,
	})
	if err != nil {
		logc.Errorf(l.ctx, "查询菜单树失败: %v", err)
		return nil, err
	}

	menuTree := logic.ConvertToTypesMenus(res.Menus)
	return &types.MenuTreeResponse{Menus: menuTree}, nil
}
