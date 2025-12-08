package authservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	bcryptUtil "zero-admin/pkg/bcrypt"
	"zero-admin/pkg/utils"
	"zero-admin/rpc/sys/db/mysql/model"

	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 注册
func (l *RegisterLogic) Register(in *sysclient.RegisterRequest) (*sysclient.RegisterResponse, error) {
	u, _ := l.svcCtx.DB.GetUserByUsername(l.ctx, in.Username)
	if u != nil && u.ID > 0 {
		logc.Errorf(l.ctx, "用户已存在, 参数：%+v", in)
		return nil, errors.New("用户已存在")
	}

	if !bcryptUtil.ValidatePasswordLength(in.Password) {
		return nil, errors.New("密码长度不符合要求")
	}
	now := time.Now()
	// 生成盐值
	salt := utils.GetRandomString(16)
	_, err := l.svcCtx.DB.CreateUser(l.ctx, model.SysUser{
		Username:      in.Username,
		Password:      bcryptUtil.HashPassword(in.Password + salt),
		Salt:          salt,
		Email:         in.Email,
		Mobile:        in.Mobile,
		RealName:      in.RealName,
		Gender:        in.Gender,
		LoginCount:    1,
		LastLoginTime: &now,
		LastLoginIP:   in.IpAddress,
		Creator:       "",
		Updater:       "",
	})
	if err != nil {
		logc.Errorf(l.ctx, "创建用户失败, 参数：%+v, 错误：%s", in, err.Error())
		return nil, status.Error(codes.Internal, "创建用户失败")
	}

	loginResp, err := NewLoginLogic(l.ctx, l.svcCtx).Login(&sysclient.LoginRequest{
		Username:  in.Username,
		Password:  in.Password,
		IpAddress: in.IpAddress,
		Os:        in.Os,
		Browser:   in.Browser,
	})
	if err != nil {
		return nil, err
	}
	return &sysclient.RegisterResponse{
		Id:           loginResp.Id,
		Username:     loginResp.Username,
		Token:        loginResp.Token,
		RefreshToken: loginResp.RefreshToken,
	}, nil
}
