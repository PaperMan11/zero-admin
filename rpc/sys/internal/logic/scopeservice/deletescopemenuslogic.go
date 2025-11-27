package scopeservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteScopeMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteScopeMenusLogic {
	return &DeleteScopeMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteScopeMenusLogic) DeleteScopeMenus(in *sysclient.DeleteScopeMenusRequest) (*sysclient.ScopeInfo, error) {
	// todo: add your logic here and delete this line

	return &sysclient.ScopeInfo{}, nil
}
