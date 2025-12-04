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

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleRequest) (resp *types.Role, err error) {
	uid := logic.GetOperateID(l.ctx)
	res, err := l.svcCtx.RoleService.CreateRole(l.ctx, &roleservice.CreateRoleRequest{
		RoleName:    req.RoleName,
		RoleCode:    req.RoleCode,
		Description: req.Description,
		Status:      req.Status,
		OperatorId:  uid,
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建角色失败: %v", err)
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
