// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/logic"
	"zero-admin/rpc/sys/sysclient"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleRoleStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleRoleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleRoleStatusLogic {
	return &ToggleRoleStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleRoleStatusLogic) ToggleRoleStatus(req *types.ToggleRoleStatusRequest) (resp *types.Role, err error) {
	uid := logic.GetOperateID(l.ctx)
	res, err := l.svcCtx.RoleService.ToggleRoleStatus(l.ctx, &sysclient.ToggleRoleStatusRequest{
		RoleId:     req.RoleId,
		Status:     req.Status,
		OperatorId: uid,
	})
	if err != nil {
		logc.Errorf(l.ctx, "切换角色状态失败: %v", err)
		return nil, err
	}

	return &types.Role{
		RoleId:      res.RoleId,
		RoleName:    res.RoleName,
		RoleCode:    res.RoleCode,
		Description: res.Description,
		Status:      res.Status,
	}, nil
}
