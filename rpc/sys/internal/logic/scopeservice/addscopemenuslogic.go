package scopeservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddScopeMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddScopeMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddScopeMenusLogic {
	return &AddScopeMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddScopeMenusLogic) AddScopeMenus(in *sysclient.AddScopeMenusRequest) (*sysclient.ScopeInfo, error) {
	// todo: add your logic here and delete this line

	return &sysclient.ScopeInfo{}, nil
}
