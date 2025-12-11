// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/authservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrRefreshTokenError   = errors.New("刷新token异常")
	ErrRefreshTokenExpired = errors.New("刷新token已过期")
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest) (resp *types.RefreshTokenResponse, err error) {
	if req.RefreshToken == "" {
		return nil, ErrRefreshTokenError
	}
	key := utils.GetRefreshTokenKey(req.Uid)
	cacheToken, _ := l.svcCtx.LocalCache.Take(key, func() (any, error) {
		return l.svcCtx.Redis.GetCtx(l.ctx, key)
	})
	if cacheToken != req.RefreshToken {
		logc.Infof(l.ctx, "cacheToken=%s, token=%s", cacheToken, req.RefreshToken)
		return nil, ErrRefreshTokenExpired
	}

	res, err := l.svcCtx.AuthService.RefreshToken(l.ctx, &authservice.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		logc.Errorf(l.ctx, "刷新token：%+v,异常:%s", req, err.Error())
		return nil, err
	}

	// 添加token过期管理
	SetTokenCache(l.ctx, l.svcCtx, res.Id, res.TokenUuid, res.RefreshToken)

	return &types.RefreshTokenResponse{
		RefreshToken: res.RefreshToken,
		AccessToken:  res.Token,
	}, nil
}
