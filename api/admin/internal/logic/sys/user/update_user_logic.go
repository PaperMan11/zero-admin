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

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (resp *types.User, err error) {
	res, err := l.svcCtx.UserService.UpdateUser(l.ctx, &userservice.UpdateUserRequest{
		Id:         req.Id,
		Email:      req.Email,
		Mobile:     req.Mobile,
		RealName:   req.RealName,
		Gender:     req.Gender,
		Avatar:     req.Avatar,
		OperatorId: logic.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "更新用户信息, 参数: %v, 异常: %v", req, err)
		return nil, err
	}

	user := logic.ConvertToTypesUser(res)
	return &user, nil
}
