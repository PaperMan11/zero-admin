// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/scopeservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateScopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScopeLogic {
	return &UpdateScopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateScopeLogic) UpdateScope(req *types.UpdateScopeRequest) (resp *types.Scope, err error) {
	res, err := l.svcCtx.ScopeService.UpdateScope(l.ctx, &scopeservice.UpdateScopeRequest{
		Id:          req.Id,
		ScopeName:   req.ScopeName,
		ScopeCode:   req.ScopeCode,
		Description: req.Description,
		Sort:        req.Sort,
		OperatorId:  utils.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "更新安全范围失败: %v", err)
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
