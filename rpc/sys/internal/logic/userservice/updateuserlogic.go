package userservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/internal/logic"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *sysclient.UpdateUserRequest) (*sysclient.User, error) {
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
		}
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 错误：%v", in, err)
		return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
	}
	user.RealName = in.RealName
	user.Mobile = in.Mobile
	user.Email = in.Email
	user.Avatar = in.Avatar
	user.Gender = in.Gender
	user.Updater = convert.ToString(in.OperatorId)
	err = l.svcCtx.DB.SaveUser(l.ctx, *user)
	if err != nil {
		logc.Errorf(l.ctx, "更新用户信息, 参数：%+v, 错误：%v", in, err)
		return nil, xerr.NewErrCodeMsg(xerr.ErrorDb, "更新用户信息失败")
	}
	return logic.ConvertToRpcUser(user), nil
}
