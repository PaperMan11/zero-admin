package userservicelogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	bcryptUtil "zero-admin/pkg/bcrypt"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/response/xerr"
	"zero-admin/pkg/utils"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *sysclient.CreateUserRequest) (*sysclient.UserInfo, error) {
	operator := convert.ToString(in.GetOperatorId())
	u, _ := l.svcCtx.DB.GetUserByUsername(l.ctx, in.Username)
	if u.ID > 0 {
		logc.Errorf(l.ctx, "用户已存在, 参数：%+v", in)
		return nil, xerr.NewErrCode(xerr.ErrorUserExist)
	}

	if !bcryptUtil.ValidatePasswordLength(in.Password) {
		return nil, xerr.NewErrCode(xerr.ErrorPasswordLength)
	}
	// 生成盐值
	salt := utils.GetRandomString(16)
	userID, err := l.svcCtx.DB.CreateUser(l.ctx, model.SysUser{
		Username: in.Username,
		Password: bcryptUtil.HashPassword(in.Password + salt),
		Salt:     salt,
		Email:    in.Email,
		Mobile:   in.Mobile,
		RealName: in.RealName,
		Gender:   in.Gender,
		Status:   in.Status,
		Creator:  operator,
		Updater:  operator,
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建用户失败, 参数：%+v, 错误：%s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorCreateUser)
	}
	user, _ := NewGetUserInfoLogic(l.ctx, l.svcCtx).GetUserInfo(&sysclient.GetUserInfoRequest{UserId: userID})
	return user, nil
}
