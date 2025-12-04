// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package operatelog

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/rpc/sys/client/operatelogservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOperateLogDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryOperateLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOperateLogDetailLogic {
	return &QueryOperateLogDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryOperateLogDetailLogic) QueryOperateLogDetail(req *types.QueryOperateLogDetailReq) (resp *types.OperateLog, err error) {
	res, err := l.svcCtx.OperateLogService.QueryOperateLogDetail(l.ctx, &operatelogservice.QueryOperateLogDetailReq{Id: req.Id})
	if err != nil {
		logc.Errorf(l.ctx, "查询系统操作日志表详情失败: %v", err)
		return nil, err
	}
	return &types.OperateLog{
		Id:                res.Id,
		Title:             res.Title,
		OperationType:     res.OperationType,
		OperationName:     res.OperationName,
		RequestMethod:     res.RequestMethod,
		OperationUrl:      res.OperationUrl,
		OperationParams:   res.OperationParams,
		OperationResponse: res.OperationResponse,
		OperationStatus:   res.OperationStatus,
		UseTime:           res.UseTime,
		Browser:           res.Browser,
		Os:                res.Os,
		OperationIp:       res.OperationIp,
		OperationTime:     res.OperationTime,
	}, nil
}
