package scopeservicelogic

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/rpc/sys/db/mysql/model"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignScopeMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignScopeMenusLogic {
	return &AssignScopeMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignScopeMenusLogic) AssignScopeMenus(in *sysclient.AssignScopeMenusRequest) (*sysclient.ScopeInfo, error) {
	scope, err := l.svcCtx.DB.GetScopeByID(l.ctx, in.ScopeId)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, status.Error(codes.NotFound, "安全范围不存在")
	case err != nil:
		return nil, status.Error(codes.Internal, "查询安全范围失败")
	}

	menus := make([]*model.SysMenu, 0)
	for _, menu := range in.Menus {
		menus = append(menus, walkMenus(menu, 0)...)
	}

	err = l.svcCtx.DB.UpdateScopeMenusTx(l.ctx, scope.ID, menus)
	if err != nil {
		return nil, status.Error(codes.Internal, "更新安全范围菜单失败")
	}

	return NewGetScopeMenusLogic(l.ctx, l.svcCtx).GetScopeMenus(&sysclient.Int64Value{Value: scope.ID})
}

// walkMenus 递归遍历菜单树并构建菜单列表
func walkMenus(menu *sysclient.AssignScopeMenuMeta, parentID int64) []*model.SysMenu {
	// 创建当前节点菜单
	menus := []*model.SysMenu{{
		ID:       menu.MenuId,
		ParentID: parentID,
	}}

	// 递归处理子菜单
	if menu.Children != nil {
		for _, child := range menu.Children {
			childMenus := walkMenus(child, menu.MenuId)
			menus = append(menus, childMenus...)
		}
	}

	return menus
}
