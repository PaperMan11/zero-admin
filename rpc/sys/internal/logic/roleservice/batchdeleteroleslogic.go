package roleservicelogic

import (
	"context"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteRolesLogic {
	return &BatchDeleteRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchDeleteRolesLogic) BatchDeleteRoles(in *sysclient.BatchDeleteRolesRequest) (*sysclient.Empty, error) {

	delRoleSvc := NewDeleteRoleLogic(l.ctx, l.svcCtx)
	for _, roleId := range in.RoleIds {
		_, err := delRoleSvc.DeleteRole(&sysclient.DeleteRoleRequest{
			Id: roleId,
		})
		if err != nil {
			return nil, err
		}
	}

	return &sysclient.Empty{}, nil
}
