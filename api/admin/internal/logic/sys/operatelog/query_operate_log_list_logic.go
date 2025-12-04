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

type QueryOperateLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryOperateLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOperateLogListLogic {
	return &QueryOperateLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryOperateLogListLogic) QueryOperateLogList(req *types.QueryOperateLogListReq) (resp *types.OperateLogListData, err error) {
	list, err := l.svcCtx.OperateLogService.QueryOperateLogList(l.ctx, &operatelogservice.QueryOperateLogListReq{
		PageRequest: &operatelogservice.PageRequest{
			Page:     int32(req.Page),
			PageSize: int32(req.PageSize),
			Keyword:  req.Keyword,
		},
		OperationIp:     req.OperationIp,
		OperationName:   req.OperationName,
		OperationStatus: req.OperationStatus,
		OperationType:   req.OperationType,
		OperationUrl:    req.OperationUrl,
		Title:           req.Title,
		Os:              req.Os,
		Browser:         req.Browser,
		RequestMethod:   req.RequestMethod,
	})
	if err != nil {
		logc.Errorf(l.ctx, "查询系统操作日志列表失败: %v", err)
		return nil, err
	}

	data := make([]types.OperateLog, 0, len(list.Data))
	for _, v := range list.Data {
		data = append(data, types.OperateLog{
			Id:                v.Id,
			Title:             v.Title,
			OperationType:     v.OperationType,
			OperationName:     v.OperationName,
			RequestMethod:     v.RequestMethod,
			OperationUrl:      v.OperationUrl,
			OperationParams:   v.OperationParams,
			OperationResponse: v.OperationResponse,
			OperationStatus:   v.OperationStatus,
			UseTime:           v.UseTime,
			Browser:           v.Browser,
			Os:                v.Os,
			OperationIp:       v.OperationIp,
			OperationTime:     v.OperationTime,
		})
	}
	return &types.OperateLogListData{
		PageResponse: types.PageResponse{
			Total:     int64(list.PageResponse.Total),
			Page:      int64(list.PageResponse.Page),
			PageSize:  int64(list.PageResponse.PageSize),
			TotalPage: int64(list.PageResponse.TotalPage),
		},
		Data: data,
	}, nil
}
