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

type GetScopeMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeMenusLogic {
	return &GetScopeMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScopeMenusLogic) GetScopeMenus(req *types.IdValue) (resp *types.ScopeInfo, err error) {
	res, err := l.svcCtx.ScopeService.GetScopeMenus(l.ctx, &scopeservice.Int64Value{Value: req.Id})
	if err != nil {
		logc.Errorf(l.ctx, "获得安全范围菜单权限失败: %v", err)
		return nil, err
	}

	return &types.ScopeInfo{
		Scope: logic.ConvertToTypesScope(res.Scope),
		Menus: logic.ConvertToTypesMenus(res.Menus),
	}, nil
}
