// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/scopeservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllScopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllScopeLogic {
	return &GetAllScopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllScopeLogic) GetAllScope(req *types.GetAllScopeRequest) (resp *types.GetAllScopeResponse, err error) {
	scopes, _ := l.svcCtx.ScopeService.GetAllScopes(l.ctx, &scopeservice.Empty{})
	return &types.GetAllScopeResponse{
		Scopes: utils.ConvertToTypesScopes(scopes.Scopes),
	}, nil
}
