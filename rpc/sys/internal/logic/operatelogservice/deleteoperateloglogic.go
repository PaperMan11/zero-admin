package operatelogservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOperateLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOperateLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOperateLogLogic {
	return &DeleteOperateLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除系统操作日志表
func (l *DeleteOperateLogLogic) DeleteOperateLog(in *sysclient.DeleteOperateLogReq) (*sysclient.DeleteOperateLogResp, error) {
	err := l.svcCtx.DB.DeleteOperateLogs(l.ctx, in.Ids)
	if err != nil {
		logc.Errorf(l.ctx, "删除系统操作日志表, 错误：%v", err)
		return &sysclient.DeleteOperateLogResp{Pong: "1"}, xerr.NewErrCodeMsg(xerr.ErrorDb, "删除系统操作日志表失败")
	}
	return &sysclient.DeleteOperateLogResp{Pong: "0"}, nil
}
