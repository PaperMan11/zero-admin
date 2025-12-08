// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/sysclient"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleRequest) (resp *types.Role, err error) {
	uid := utils.GetOperateID(l.ctx)
	res, err := l.svcCtx.RoleService.UpdateRole(l.ctx, &sysclient.UpdateRoleRequest{
		RoleName:    req.RoleName,
		RoleCode:    req.RoleCode,
		Description: req.Description,
		Status:      req.Status,
		OperatorId:  uid,
	})
	if err != nil {
		logc.Errorf(l.ctx, "更新角色失败: %v", err)
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
