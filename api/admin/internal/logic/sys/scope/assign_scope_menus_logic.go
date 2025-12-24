// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/scopeservice"
)

type AssignScopeMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignScopeMenusLogic {
	return &AssignScopeMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignScopeMenusLogic) AssignScopeMenus(req *types.AssignScopeMenusRequest) (resp *types.ScopeInfo, err error) {
	menus := make([]*scopeservice.AssignScopeMenuMeta, 0)
	for _, menu := range req.Menus {
		menus = append(menus, convertToRpcAssignScopeMenuMeta(&menu))
	}
	scopeMenus, err := l.svcCtx.ScopeService.AssignScopeMenus(l.ctx, &scopeservice.AssignScopeMenusRequest{
		ScopeId:    req.ScopeId,
		Menus:      menus,
		OperatorId: utils.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "分配菜单失败，安全范围：%d，错误：%v", req.ScopeId, err)
		return
	}

	return &types.ScopeInfo{
		Scope: utils.ConvertToTypesScope(scopeMenus.Scope),
		Menus: utils.ConvertToTypesMenus(scopeMenus.Menus),
	}, nil
}

func convertToRpcAssignScopeMenuMeta(menu *types.Menu) *scopeservice.AssignScopeMenuMeta {
	res := &scopeservice.AssignScopeMenuMeta{
		MenuId:   menu.Id,
		ParentId: menu.ParentId,
		Children: make([]*scopeservice.AssignScopeMenuMeta, 0),
	}
	for _, child := range menu.Children {
		res.Children = append(res.Children, convertToRpcAssignScopeMenuMeta(&child))
	}
	return res
}
