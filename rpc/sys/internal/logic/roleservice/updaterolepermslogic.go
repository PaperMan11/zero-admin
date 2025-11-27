package roleservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRolePermsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRolePermsLogic {
	return &UpdateRolePermsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色权限
func (l *UpdateRolePermsLogic) UpdateRolePerms(in *sysclient.UpdateRolePermsRequest) (*sysclient.RoleInfo, error) {
	// todo: add your logic here and delete this line

	return &sysclient.RoleInfo{}, nil
}
