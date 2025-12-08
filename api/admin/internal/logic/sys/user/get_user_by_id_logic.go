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

type GetUserByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserByIdLogic) GetUserById(req *types.IdValue) (resp *types.UserInfo, err error) {
	res, err := l.svcCtx.UserService.GetUserById(l.ctx, &userservice.Int64Value{Value: req.Id})
	if err != nil {
		logc.Errorf(l.ctx, "获取用户信息失败: %v", err)
		return nil, err
	}

	userInfo := utils.ConvertToTypesUserInfo(res)
	return &userInfo, nil
}
