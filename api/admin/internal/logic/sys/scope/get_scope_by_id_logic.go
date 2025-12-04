// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScopeByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetScopeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeByIdLogic {
	return &GetScopeByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScopeByIdLogic) GetScopeById(req *types.IdValue) (resp *types.Scope, err error) {
	// todo: add your logic here and delete this line

	return
}
