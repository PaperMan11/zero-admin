package operatelogservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/logic"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOperateLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryOperateLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOperateLogListLogic {
	return &QueryOperateLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询系统操作日志表列表
func (l *QueryOperateLogListLogic) QueryOperateLogList(in *sysclient.QueryOperateLogListReq) (*sysclient.OperateLogListData, error) {
	logs, total, err := l.svcCtx.DB.GetOperateLogs(l.ctx, model.OperateLogFilter{
		Title:           in.Title,
		OperationType:   in.OperationType,
		OperationName:   in.OperationName,
		RequestMethod:   in.RequestMethod,
		OperationURL:    in.OperationUrl,
		OperationStatus: in.OperationStatus,
		Browser:         in.Browser,
		Os:              in.Os,
		OperationIP:     in.OperationIp,
	}, int(in.PageRequest.Page), int(in.PageRequest.PageSize))
	if err != nil {
		logc.Errorf(l.ctx, "查询系统操作日志表列表失败: %v", err)
		return nil, xerr.NewErrCodeMsg(xerr.ErrorDb, "查询系统操作日志表列表失败")
	}
	return &sysclient.OperateLogListData{
		PageResponse: &sysclient.PageResponse{
			Total:     int32(total),
			Page:      in.PageRequest.Page,
			PageSize:  in.PageRequest.PageSize,
			TotalPage: int32(total)/in.PageRequest.PageSize + 1,
		},
		Data: logic.ConvertToRpcOperateLogs(logs),
	}, nil
}
