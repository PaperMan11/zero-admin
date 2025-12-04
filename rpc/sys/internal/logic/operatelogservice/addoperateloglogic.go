package operatelogservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/rpc/sys/db/mysql/model"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddOperateLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddOperateLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddOperateLogLogic {
	return &AddOperateLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加系统操作日志表
func (l *AddOperateLogLogic) AddOperateLog(in *sysclient.AddOperateLogReq) (*sysclient.AddOperateLogResp, error) {
	_, err := l.svcCtx.DB.CreateOperationLog(l.ctx, model.SysOperateLog{
		Title:             in.Title,
		OperationType:     in.OperationType,
		OperationName:     in.OperationName,
		RequestMethod:     in.RequestMethod,
		OperationURL:      in.OperationUrl,
		OperationParams:   in.OperationParams,
		OperationResponse: in.OperationResponse,
		OperationStatus:   in.OperationStatus,
		UseTime:           in.UseTime,
		Browser:           in.Browser,
		Os:                in.Os,
		OperationIP:       in.OperationIp,
	})
	if err != nil {
		logc.Errorf(l.ctx, "保存系统操作日志异常, 登录参数：%+v, 错误：%s", in, err.Error())
		return &sysclient.AddOperateLogResp{Pong: "1"}, status.Error(codes.Internal, "保存系统操作日志异常")
	}
	return &sysclient.AddOperateLogResp{Pong: "0"}, nil
}
