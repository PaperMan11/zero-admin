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

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserRequest) (resp *types.Empty, err error) {
	_, err = l.svcCtx.UserService.DeleteUser(l.ctx, &userservice.DeleteUserRequest{
		Id:         req.Id,
		OperatorId: logic.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "删除用户失败: %v", err)
		return nil, err
	}
	return &types.Empty{}, nil
}
