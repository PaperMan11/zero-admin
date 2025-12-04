// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/logic"
	"zero-admin/rpc/sys/client/userservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignUserRoleLogic {
	return &AssignUserRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignUserRoleLogic) AssignUserRole(req *types.AssignUserRoleRequest) (resp *types.Empty, err error) {
	_, err = l.svcCtx.UserService.AssignUserRole(l.ctx, &userservice.AssignUserRoleRequest{
		UserId:     req.UserId,
		RoleCodes:  req.RoleCodes,
		OperatorId: logic.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "分配用户角色, 参数: %v, 异常: %v", req, err)
		return nil, err
	}

	return &types.Empty{}, nil
}
