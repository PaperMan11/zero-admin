package roleservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolePermsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermsLogic {
	return &GetRolePermsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取角色权限
func (l *GetRolePermsLogic) GetRolePerms(in *sysclient.Int64Value) (*sysclient.RoleInfo, error) {
	// todo: add your logic here and delete this line

	return &sysclient.RoleInfo{}, nil
}
