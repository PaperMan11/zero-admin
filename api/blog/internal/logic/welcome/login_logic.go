// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package welcome

import (
	"context"
	"errors"
	"time"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	bcryptUtil "zero-admin/pkg/bcrypt"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/errorx"
	jwtUtil "zero-admin/pkg/jwt"
	"zero-admin/pkg/utils"
	uuidUtil "zero-admin/pkg/uuid"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq, clientIP string) (resp *types.LoginResp, err error) {
	if !l.svcCtx.Captcha.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, errorx.ErrorCaptcha
	}

	user, err := l.svcCtx.UserModel.GetByUsername(l.ctx, req.Username)
	// 1.判断用户是否存在
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		l.Errorf("用户不存在, 参数：%+v, 异常: %s", req, err.Error())
		return nil, errorx.ErrorInvalidUsernameOrPassword
	case err != nil:
		l.Errorf("查询用户信息, 参数：%+v, 异常: %s", req, err.Error())
		return nil, errorx.ErrorInternal
	}
	if user.Status != 0 {
		return nil, errorx.ErrorUserDisabled
	}
	if user.DeletedAt.Valid {
		return nil, errorx.ErrorUserNotFound
	}

	// 2.判断密码是否正确
	if !bcryptUtil.CheckPassword(req.Password+user.Salt, user.Password) {
		return nil, errorx.ErrorInvalidUsernameOrPassword
	}

	// 登录成功后，返回token
	uuid, accessToken, refreshToken, err := GenerateToken(int64(user.ID), []string{"user"},
		l.svcCtx.Config.Name, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire,
		l.svcCtx.Config.Auth.RefreshSecret, l.svcCtx.Config.Auth.RefreshExpire)
	if err != nil {
		return nil, err
	}

	l.svcCtx.UserModel.Updates(l.ctx, user.ID, map[string]interface{}{
		"last_login_ip":   clientIP,
		"last_login_time": time.Now(),
		"login_count":     user.LoginCount + 1,
	})

	// 3.缓存token
	keyAccessToken := utils.GetAccessTokenKey(int64(user.ID), uuid)
	keyRefreshToken := utils.GetRefreshTokenKey(int64(user.ID), uuid)
	accessExpiredAt := time.Now().Add(time.Second * time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Unix()
	refreshExpiredAt := time.Now().Add(time.Second * time.Duration(l.svcCtx.Config.Auth.RefreshExpire)).Unix()
	l.svcCtx.LocalCache.SetWithExpire(keyAccessToken, accessExpiredAt, time.Second*time.Duration(l.svcCtx.Config.Auth.AccessExpire))
	l.svcCtx.LocalCache.SetWithExpire(keyRefreshToken, refreshExpiredAt, time.Second*time.Duration(l.svcCtx.Config.Auth.RefreshExpire))
	l.svcCtx.Redis.SetexCtx(l.ctx, keyAccessToken, convert.ToString(accessExpiredAt), int(l.svcCtx.Config.Auth.AccessExpire))
	l.svcCtx.Redis.SetexCtx(l.ctx, keyRefreshToken, convert.ToString(refreshExpiredAt), int(l.svcCtx.Config.Auth.RefreshExpire))

	return &types.LoginResp{
		UserId:       int64(user.ID),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func GenerateToken(userID int64, roles []string,
	issuer string, accessSecret string, accessExpire int64,
	refreshSecret string, refreshExpire int64) (uuid, accessToken, refreshToken string, err error) {
	uuid = uuidUtil.GetUUID()
	accessToken, err = jwtUtil.GenerateAccessToken(uuid, issuer, userID, roles, accessSecret, accessExpire)
	if err != nil {
		return "", "", "", errorx.ErrorTokenGenerate
	}
	refreshToken, err = jwtUtil.GenerateRefreshToken(uuid, issuer, userID, roles, refreshSecret, refreshExpire)
	if err != nil {
		return "", "", "", errorx.ErrorTokenGenerate
	}
	return uuid, accessToken, refreshToken, nil
}
