// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package welcome

import (
	"context"
	"time"

	"zero-admin/api/blog/internal/models"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/bcrypt"
	bcryptUtil "zero-admin/pkg/bcrypt"
	"zero-admin/pkg/convert"
	"zero-admin/pkg/errorx"
	"zero-admin/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq, clientIP string) (resp *types.RegisterResp, err error) {
	if req.Email != "" {
		user, _ := l.svcCtx.UserModel.GetByEmail(l.ctx, req.Email)
		if user.Email == req.Email {
			return nil, errorx.ErrorEmailRegistered
		}
	}
	if req.Mobile != "" {
		user, _ := l.svcCtx.UserModel.GetByMobile(l.ctx, req.Mobile)
		if user.Mobile == req.Mobile {
			return nil, errorx.ErrorMobileRegistered
		}
	}
	if req.Username != "" {
		user, _ := l.svcCtx.UserModel.GetByUsername(l.ctx, req.Username)
		if user.Username == req.Username {
			return nil, errorx.ErrorUsernameRegistered
		}
	}
	if !bcryptUtil.ValidatePasswordLength(req.Password) {
		return nil, errorx.ErrorInvalidPasswordLen
	}
	// 生成盐值
	salt := utils.GetRandomString(16)
	bcryptPassword := bcrypt.HashPassword(req.Password + salt)

	// 注册用户
	now := time.Now()
	user := models.User{
		Username:     req.Username,
		Nickname:     req.Nickname,
		Password:     bcryptPassword,
		Salt:         salt,
		Avatar:       req.Avatar,       // 头像
		Gender:       int8(req.Gender), // 性别 0-未知 1-男 2-女
		Introduction: req.Introduction, // 个人介绍
		Birthday:     req.Birthday,     // 出生日期
		Mobile:       req.Mobile,       // 手机号
		Email:        req.Email,        // 邮箱
		LoginTime:    now,              // 最近登录时间
		LastLoginIP:  clientIP,         // 最近登录IP
		LoginCount:   1,                // 登录次数
	}
	userId, err := l.svcCtx.UserModel.Create(l.ctx, user)
	if err != nil {
		l.Errorf("注册用户失败: %v", err)
		return nil, errorx.ErrorInternal
	}

	// 登录成功后，返回token
	uuid, accessToken, refreshToken, err := GenerateToken(int64(userId), []string{"user"},
		l.svcCtx.Config.Name, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire,
		l.svcCtx.Config.Auth.RefreshSecret, l.svcCtx.Config.Auth.RefreshExpire)
	if err != nil {
		l.Errorf("生成token失败: %v", err)
		return nil, errorx.ErrorInternal
	}

	// 3.缓存token
	keyAccessToken := utils.GetAccessTokenKey(int64(userId), uuid)
	keyRefreshToken := utils.GetRefreshTokenKey(int64(userId), uuid)
	accessExpiredAt := time.Now().Add(time.Second * time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Unix()
	refreshExpiredAt := time.Now().Add(time.Second * time.Duration(l.svcCtx.Config.Auth.RefreshExpire)).Unix()
	l.svcCtx.LocalCache.SetWithExpire(keyAccessToken, accessExpiredAt, time.Second*time.Duration(l.svcCtx.Config.Auth.AccessExpire))
	l.svcCtx.LocalCache.SetWithExpire(keyRefreshToken, refreshExpiredAt, time.Second*time.Duration(l.svcCtx.Config.Auth.RefreshExpire))
	l.svcCtx.Redis.SetexCtx(l.ctx, keyAccessToken, convert.ToString(accessExpiredAt), int(l.svcCtx.Config.Auth.AccessExpire))
	l.svcCtx.Redis.SetexCtx(l.ctx, keyRefreshToken, convert.ToString(refreshExpiredAt), int(l.svcCtx.Config.Auth.RefreshExpire))

	return &types.RegisterResp{
		UserId:       int64(userId),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
