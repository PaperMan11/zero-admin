// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/userservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(req *types.UpdatePasswordRequest) (resp *types.Empty, err error) {
	_, err = l.svcCtx.UserService.UpdatePassword(l.ctx, &userservice.UpdatePasswordRequest{
		UserId:     req.UserId,
		Password:   req.Password,
		OperatorId: utils.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "管理员更新用户密码, 参数: %v, 异常: %v", req, err)
		return nil, err
	}
	return &types.Empty{}, nil
}
