// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.UserInfo, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	user, err := l.svcCtx.UserModel.GetByID(l.ctx, uint(uid))
	if err != nil {
		l.Errorf("检查用户状态失败: %v, 用户ID: %d", err, uid)
		return nil, err
	}
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if err := l.svcCtx.UserModel.Updates(l.ctx, user.ID, updates); err != nil {
		l.Errorf("更新用户信息失败: %v, 用户ID: %d", err, uid)
		return nil, errorx.ErrorInternal
	}

	return &types.UserInfo{
		Id:       convert.ToInt64(user.ID),
		Nickname: user.Nickname,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}, nil
}
