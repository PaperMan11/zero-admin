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

type UpdateUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(req *types.UpdateUserPasswordRequest) (resp *types.Empty, err error) {
	_, err = l.svcCtx.UserService.UpdateUserPassword(l.ctx, &userservice.UpdateUserPasswordRequest{
		UserId:      req.UserId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
		OperatorId:  logic.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "更新用户密码, 参数: %v, 异常: %v", req, err)
		return nil, err
	}
	return &types.Empty{}, nil
}
