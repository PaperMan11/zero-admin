// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package operatelog

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/rpc/sys/sysclient"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOperateLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOperateLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOperateLogLogic {
	return &DeleteOperateLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOperateLogLogic) DeleteOperateLog(req *types.DeleteOperateLogReq) (resp *types.DeleteOperateLogResp, err error) {
	logResp, err := l.svcCtx.OperateLogService.DeleteOperateLog(l.ctx, &sysclient.DeleteOperateLogReq{
		Ids: req.Ids,
	})
	if err != nil {
		logc.Errorf(l.ctx, "删除系统操作日志表失败: %v", err)
		return nil, err
	}
	return &types.DeleteOperateLogResp{
		Pong: logResp.Pong,
	}, nil
}
