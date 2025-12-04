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

type ToggleUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleUserStatusLogic {
	return &ToggleUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleUserStatusLogic) ToggleUserStatus(req *types.ToggleUserStatusRequest) (resp *types.User, err error) {
	res, err := l.svcCtx.UserService.ToggleUserStatus(l.ctx, &userservice.ToggleUserStatusRequest{
		UserId:     req.UserId,
		Status:     req.Status,
		OperatorId: logic.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "修改用户启用状态失败: %v", err)
		return nil, err
	}

	user := logic.ConvertToTypesUser(res)
	return &user, nil
}
