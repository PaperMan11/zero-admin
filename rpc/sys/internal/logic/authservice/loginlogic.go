package authservicelogic

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"time"
	bcryptUtil "zero-admin/pkg/bcrypt"
	jwtUtil "zero-admin/pkg/jwt"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登录
func (l *LoginLogic) Login(in *sysclient.LoginRequest) (*sysclient.LoginResponse, error) {
	user, err := l.svcCtx.DB.GetUserByUsername(l.ctx, in.Username)
	// 1.判断用户是否存在
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		logc.Errorf(l.ctx, "用户不存在, 参数：%+v, 异常: %s", in, err.Error())
		l.saveLoginLog(in, LoginStatusFail, "用户不存在")
		return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
	case err != nil:
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		l.saveLoginLog(in, LoginStatusFail, fmt.Sprintf("系统异常:%s", err.Error()))
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	// 2.判断密码是否正确
	if !bcryptUtil.CheckPassword(in.Password, user.Password) {
		l.saveLoginLog(in, LoginStatusFail, "密码错误")
		return nil, xerr.NewErrCode(xerr.ErrorUserPassword)
	}

	// 3.生成token
	// 用户角色信息
	_, roleCodes := GetUserRoles(l.ctx, l.svcCtx.DB, user.ID)
	accessToken, refreshToken, err := GenerateToken(user.ID, roleCodes, l.svcCtx.Config.Name,
		l.svcCtx.Config.Jwt.AccessSecret, l.svcCtx.Config.Jwt.AccessExpire,
		l.svcCtx.Config.Jwt.RefreshSecret, l.svcCtx.Config.Jwt.RefreshExpire)
	if err != nil {
		logc.Errorf(l.ctx, "生成token异常, 登录参数：%+v, 错误：%s", in, err.Error())
		l.saveLoginLog(in, LoginStatusFail, "生成token异常")
		return nil, xerr.NewErrCode(xerr.ErrorTokenGenerate)
	}
	// 4.更新登录时间
	_ = l.svcCtx.DB.UpdateUserByID(l.ctx, user.ID, map[string]any{"last_login_time": time.Now()})

	l.saveLoginLog(in, LoginStatusSuccess, "")
	return &sysclient.LoginResponse{
		Id:           user.ID,
		Username:     user.Username,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

const (
	LoginStatusSuccess int32 = 1
	LoginStatusFail    int32 = 0
)

func (l *LoginLogic) saveLoginLog(in *sysclient.LoginRequest, status int32, msg string) {
	if status == LoginStatusSuccess {
		msg = "登录成功"
	}
	err := l.svcCtx.DB.CreateLoginLog(l.ctx, model.SysLoginLog{
		Username:  in.Username,
		IP:        in.IpAddress,
		Location:  "",
		Browser:   in.Browser,
		Os:        in.Os,
		Status:    status,
		Message:   msg,
		LoginTime: time.Now(),
	})
	if err != nil {
		logc.Errorf(l.ctx, "保存登录日志异常, 登录参数：%+v, 错误：%s", in, err.Error())
	}
}

func GenerateToken(userID int64, roles []string,
	issuer string, accessSecret string, accessExpire int64,
	refreshSecret string, refreshExpire int64) (accessToken, refreshToken string, err error) {
	accessToken, err = jwtUtil.GenerateAccessToken(issuer, userID, roles, accessSecret, accessExpire)
	if err != nil {
		return "", "", xerr.NewErrCode(xerr.ErrorTokenGenerate)
	}
	refreshToken, _ = jwtUtil.GenerateRefreshToken(issuer, userID, refreshSecret, refreshExpire)
	return accessToken, refreshToken, nil
}
