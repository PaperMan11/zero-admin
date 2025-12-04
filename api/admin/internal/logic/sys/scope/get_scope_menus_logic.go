// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
