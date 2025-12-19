// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	//menus := make([]*scopeservice.AssignScopeMenuMeta, 0)
	//l.svcCtx.ScopeService.AssignScopeMenus(l.ctx, req.ScopeId, req.Menus)

	return
}
