package scopeservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScopeMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeMenusLogic {
	return &GetScopeMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScopeMenusLogic) GetScopeMenus(in *sysclient.Int64Value) (*sysclient.ScopeInfo, error) {
	// todo: add your logic here and delete this line

	return &sysclient.ScopeInfo{}, nil
}
