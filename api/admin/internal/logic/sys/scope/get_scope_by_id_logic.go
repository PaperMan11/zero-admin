// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/rpc/sys/client/scopeservice"

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
	res, err := l.svcCtx.ScopeService.GetScopeById(l.ctx, &scopeservice.Int64Value{Value: req.Id})
	if err != nil {
		logc.Errorf(l.ctx, "查询安全范围信息失败: %v", err)
		return nil, err
	}

	return &types.Scope{
		Id:          res.Id,
		ScopeName:   res.ScopeName,
		ScopeCode:   res.ScopeCode,
		Description: res.Description,
		Sort:        res.Sort,
	}, nil
}
