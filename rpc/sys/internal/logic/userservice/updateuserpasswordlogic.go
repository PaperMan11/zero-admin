package userservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"zero-admin/pkg/bcrypt"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(in *sysclient.UpdateUserPasswordRequest) (*sysclient.Empty, error) {
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
		}
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	if !bcrypt.ValidatePasswordLength(in.NewPassword) {
		return nil, xerr.NewErrCode(xerr.ErrorPasswordLength)
	}

	if !bcrypt.CheckPassword(in.OldPassword+user.Salt, user.Password) {
		return nil, xerr.NewErrMsg("旧密码错误")
	}
	if bcrypt.CheckPassword(in.NewPassword+user.Salt, user.Password) {
		return nil, xerr.NewErrMsg("新密码不能与旧密码相同")
	}

	hashPass := bcrypt.HashPassword(in.NewPassword + user.Salt)
	err = l.svcCtx.DB.UpdateUserByID(l.ctx, in.UserId, map[string]interface{}{"password": hashPass, "updater": convert.ToString(in.OperatorId)})
	if err != nil {
		logc.Errorf(l.ctx, "更新用户密码, 错误：%s", err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	return &sysclient.Empty{}, nil
}
