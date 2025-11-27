package scopeservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateScopeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScopeLogic {
	return &UpdateScopeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateScopeLogic) UpdateScope(in *sysclient.UpdateScopeRequest) (*sysclient.Scope, error) {
	// todo: add your logic here and delete this line

	return &sysclient.Scope{}, nil
}
