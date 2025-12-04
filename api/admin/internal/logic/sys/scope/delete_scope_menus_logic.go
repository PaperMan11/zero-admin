// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScopeMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScopeMenusLogic {
	return &DeleteScopeMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteScopeMenusLogic) DeleteScopeMenus(req *types.DeleteScopeMenusRequest) (resp *types.ScopeInfo, err error) {
	// todo: add your logic here and delete this line

	return
}
