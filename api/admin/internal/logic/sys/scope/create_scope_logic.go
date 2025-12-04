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

type CreateScopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateScopeLogic {
	return &CreateScopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateScopeLogic) CreateScope(req *types.CreateScopeRequest) (resp *types.Scope, err error) {
	res, err := l.svcCtx.ScopeService.CreateScope(l.ctx, &scopeservice.CreateScopeRequest{
		ScopeName:   req.ScopeName,
		ScopeCode:   req.ScopeCode,
		Description: req.Description,
		Sort:        req.Sort,
		MenuIds:     req.MenuIds,
		OperatorId:  logic.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建安全范围失败: %v", err)
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
