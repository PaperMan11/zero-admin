// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/logic"
	"zero-admin/rpc/sys/client/roleservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.DeleteRoleRequest) (resp *types.Empty, err error) {
	uid := logic.GetOperateID(l.ctx)
	_, err = l.svcCtx.RoleService.DeleteRole(l.ctx, &roleservice.DeleteRoleRequest{
		Id:         req.Id,
		OperatorId: uid,
	})
	if err != nil {
		logc.Errorf(l.ctx, "删除角色失败: %v", err)
		return nil, err
	}
	return &types.Empty{}, nil
}
