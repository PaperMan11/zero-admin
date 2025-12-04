// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScopeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetScopeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeListLogic {
	return &GetScopeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScopeListLogic) GetScopeList(req *types.ScopeListRequest) (resp *types.ScopeListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
