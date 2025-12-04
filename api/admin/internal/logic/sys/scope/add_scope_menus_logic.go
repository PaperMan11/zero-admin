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

type AddScopeMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddScopeMenusLogic {
	return &AddScopeMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddScopeMenusLogic) AddScopeMenus(req *types.AddScopeMenusRequest) (resp *types.ScopeInfo, err error) {
	res, err := l.svcCtx.ScopeService.AddScopeMenus(l.ctx, &scopeservice.AddScopeMenusRequest{
		ScopeId:    req.ScopeId,
		MenuIds:    req.MenuIds,
		OperatorId: logic.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "添加安全范围权限失败: %v", err)
		return nil, err
	}

	return &types.ScopeInfo{
		Scope: types.Scope{
			Id:          res.Scope.Id,
			ScopeName:   res.Scope.ScopeName,
			ScopeCode:   res.Scope.ScopeCode,
			Description: res.Scope.Description,
			Sort:        res.Scope.Sort,
		},
		Menus: ConvertToTypesMenus(res.Menus),
	}, nil
}

func ConvertToTypesMenus(menus []*scopeservice.Menu) []types.Menu {
	menuList := make([]types.Menu, 0, len(menus))
	for _, v := range menus {
		menuList = append(menuList, ConvertToTypesMenu(v))
	}
	return menuList
}

func ConvertToTypesMenu(menu *scopeservice.Menu) types.Menu {
	return types.Menu{
		Id:        menu.Id,
		ParentId:  menu.ParentId,
		MenuType:  menu.MenuType,
		Path:      menu.Path,
		Component: menu.Component,
		Redirect:  menu.Redirect,
		Sort:      menu.Sort,
		External:  menu.External,
		Hidden:    menu.Hidden,
		Status:    menu.Status,
		ScopeId:   menu.ScopeId,
		Remark:    menu.Remark,
		Meta: types.MenuMeta{
			MenuName: menu.MenuName,
			Icon:     menu.Icon,
			NoCache:  menu.NoCache,
			Affix:    menu.Affix,
		},
		Children: ConvertToTypesMenus(menu.Children),
		Perms:    menu.Perms,
	}
}
