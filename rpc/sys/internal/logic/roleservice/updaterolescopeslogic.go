package roleservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleScopesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleScopesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleScopesLogic {
	return &UpdateRoleScopesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleScopesLogic) UpdateRoleScopes(in *sysclient.UpdateRoleScopesRequest) (*sysclient.RoleInfo, error) {
	// todo: add your logic here and delete this line

	return &sysclient.RoleInfo{}, nil
}
