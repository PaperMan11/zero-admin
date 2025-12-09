// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/authservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *LoginLogic) Login(req *types.LoginRequest, ip, os, browser string) (resp *types.LoginResponse, err error) {
	res, err := l.svcCtx.AuthService.Login(l.ctx, &authservice.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		IpAddress: ip,
		Os:        os,
		Browser:   browser,
	})
	if err != nil {
		logc.Errorf(l.ctx, "用户登录：%+v,异常:%s", req, err.Error())
		return nil, err
	}

	// 添加token过期管理
	SetTokenCache(l.ctx, l.svcCtx, res.Id, res.TokenUuid, res.RefreshToken)

	return &types.LoginResponse{
		AccessToken:  res.Token,
		Id:           res.Id,
		RefreshToken: res.RefreshToken,
		Username:     res.Username,
	}, nil
}

// 添加token过期管理
func SetTokenCache(ctx context.Context, svcCtx *svc.ServiceContext, uid int64, tokenUuid, refreshToken string) {
	accessTokenExpire := svcCtx.Config.Auth.AccessExpire
	refreshTokenExpire := svcCtx.Config.Auth.RefreshExpire
	aTokenKey := utils.GetAccessTokenKey(uid)
	rTokenKey := utils.GetRefreshTokenKey(uid)
	svcCtx.LocalCache.SetWithExpire(aTokenKey, tokenUuid, time.Duration(accessTokenExpire)*time.Second)
	svcCtx.LocalCache.SetWithExpire(rTokenKey, refreshToken, time.Duration(refreshTokenExpire)*time.Second)
	svcCtx.Redis.SetexCtx(ctx, aTokenKey, tokenUuid, int(accessTokenExpire))
	svcCtx.Redis.SetexCtx(ctx, rTokenKey, refreshToken, int(refreshTokenExpire))
}

func DelTokenCache(ctx context.Context, svcCtx *svc.ServiceContext, uid int64) {
	aTokenKey := utils.GetAccessTokenKey(uid)
	rTokenKey := utils.GetRefreshTokenKey(uid)
	svcCtx.LocalCache.Del(aTokenKey)
	svcCtx.LocalCache.Del(rTokenKey)
	svcCtx.Redis.DelCtx(ctx, aTokenKey, rTokenKey)
}
