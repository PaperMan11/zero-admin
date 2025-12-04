// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/rpc/sys/client/roleservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteRolesLogic {
	return &BatchDeleteRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteRolesLogic) BatchDeleteRoles(req *types.BatchDeleteRolesRequest) (resp *types.Empty, err error) {
	_, err = l.svcCtx.RoleService.BatchDeleteRoles(l.ctx, &roleservice.BatchDeleteRolesRequest{
		OperatorId: req.OperatorId,
		RoleIds:    req.RoleIds,
	})
	if err != nil {
		logc.Errorf(l.ctx, "删除角色失败: %v", err)
		return nil, err
	}

	return &types.Empty{}, nil
}
