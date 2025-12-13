package scopeservicelogic

import (
	"context"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllScopesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllScopesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllScopesLogic {
	return &GetAllScopesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 安全范围管理
func (l *GetAllScopesLogic) GetAllScopes(in *sysclient.Empty) (*sysclient.GetAllScopesResponse, error) {
	res, _ := l.svcCtx.DB.GetAllScopes(l.ctx)
	return &sysclient.GetAllScopesResponse{
		Scopes: logic.ConvertToRpcScopes(res),
	}, nil
}
