// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/sysclient"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScopeLogic {
	return &DeleteScopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteScopeLogic) DeleteScope(req *types.DeleteScopeRequest) (resp *types.Empty, err error) {
	_, err = l.svcCtx.ScopeService.DeleteScope(l.ctx, &sysclient.DeleteScopeRequest{
		Id:         req.Id,
		OperatorId: utils.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "删除安全范围失败: %v", err)
		return nil, err
	}

	return &types.Empty{}, nil
}
