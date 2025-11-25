package permissionservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScopeByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScopeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScopeByIdLogic {
	return &GetScopeByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetScopeByIdLogic) GetScopeById(in *sysclient.Int64Value) (*sysclient.ScopeInfo, error) {
	// todo: add your logic here and delete this line

	return &sysclient.ScopeInfo{}, nil
}
