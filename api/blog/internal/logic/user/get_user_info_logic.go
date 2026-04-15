// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.Empty) (resp *types.UserInfo, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	user, err := l.svcCtx.UserModel.GetByID(l.ctx, uint(uid))
	if err != nil {
		l.Errorf("检查用户状态失败: %v, 用户ID: %d", err, uid)
		return nil, errorx.ErrorInternal
	}
	return &types.UserInfo{
		Id:           int64(user.ID),
		Username:     user.Username,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Gender:       int(user.Gender),
		Introduction: user.Introduction,
		Birthday:     user.Birthday,
		Mobile:       user.Mobile,
		Email:        user.Email,
	}, nil
}
