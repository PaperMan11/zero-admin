package roleservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRolePermsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRolePermsLogic {
	return &DeleteRolePermsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除角色权限
func (l *DeleteRolePermsLogic) DeleteRolePerms(in *sysclient.DeleteRolePermsRequest) (*sysclient.RoleInfo, error) {
	// todo: add your logic here and delete this line

	return &sysclient.RoleInfo{}, nil
}
