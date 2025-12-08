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

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserRequest) (resp *types.UserInfo, err error) {
	res, err := l.svcCtx.UserService.CreateUser(l.ctx, &userservice.CreateUserRequest{
		Username:   req.Username,
		Password:   req.Password,
		Email:      req.Email,
		Mobile:     req.Mobile,
		RealName:   req.RealName,
		Gender:     req.Gender,
		Status:     req.Status,
		RoleIds:    req.RoleIds,
		OperatorId: utils.GetOperateID(l.ctx),
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建用户失败: %v", err)
		return nil, err
	}

	userInfo := utils.ConvertToTypesUserInfo(res)
	return &userInfo, nil
}
