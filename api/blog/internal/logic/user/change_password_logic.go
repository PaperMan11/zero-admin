// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	bcryptUtil "zero-admin/pkg/bcrypt"
	"zero-admin/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.Empty, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	user, err := l.svcCtx.UserModel.GetByID(l.ctx, uint(uid))
	if err != nil {
		l.Errorf("检查用户状态失败: %v, 用户ID: %d", err, uid)
		return nil, err
	}
	if !bcryptUtil.ValidatePasswordLength(req.NewPassword) {
		return nil, errorx.ErrorInvalidPasswordLen
	}
	if !bcryptUtil.CheckPassword(req.OldPassword+user.Salt, user.Password) {
		return nil, errorx.ErrorInvalidOldPassword
	}
	user.Password = bcryptUtil.HashPassword(req.NewPassword + user.Salt)
	if err := l.svcCtx.UserModel.Updates(l.ctx, user.ID, map[string]interface{}{
		"password": user.Password,
	}); err != nil {
		l.Errorf("更新用户密码失败: %v, 用户ID: %d", err, uid)
		return &types.Empty{}, nil
	}
	return &types.Empty{}, nil
}
