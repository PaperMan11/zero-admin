package scopeservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScopeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScopeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeListLogic {
	return &GetScopeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 安全范围管理
func (l *GetScopeListLogic) GetScopeList(in *sysclient.ScopeListRequest) (*sysclient.ScopeListResponse, error) {
	// todo: add your logic here and delete this line

	return &sysclient.ScopeListResponse{}, nil
}
