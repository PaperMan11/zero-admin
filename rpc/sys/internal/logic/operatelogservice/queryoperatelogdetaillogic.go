package operatelogservicelogic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOperateLogDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryOperateLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOperateLogDetailLogic {
	return &QueryOperateLogDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询系统操作日志表详情
func (l *QueryOperateLogDetailLogic) QueryOperateLogDetail(in *sysclient.QueryOperateLogDetailReq) (*sysclient.OperateLog, error) {
	log, err := l.svcCtx.DB.GetOperateLog(l.ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询系统操作日志表详情失败")
	}
	return logic.ConvertToRpcOperateLog(log), nil
}
