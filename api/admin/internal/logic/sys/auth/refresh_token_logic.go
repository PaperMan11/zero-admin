// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/pkg/convert"
	"zero-admin/rpc/sys/client/authservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	res, err := l.svcCtx.AuthService.RefreshToken(l.ctx, &authservice.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		logc.Errorf(l.ctx, "刷新token：%+v,异常:%s", req, err.Error())
		return nil, err
	}

	// 添加token过期管理
	uid := convert.ToInt64(l.ctx.Value("uid"))
	accessTokenExpire := l.svcCtx.Config.Auth.AccessExpire
	refreshTokenExpire := l.svcCtx.Config.Auth.RefreshExpire
	l.svcCtx.Redis.SetexCtx(l.ctx, utils.GetAccessTokenKey(uid), res.Token, int(accessTokenExpire))
	l.svcCtx.Redis.SetexCtx(l.ctx, utils.GetRefreshTokenKey(uid), res.RefreshToken, int(refreshTokenExpire))

	return &types.RefreshTokenResponse{
		RefreshToken: res.RefreshToken,
		Token:        res.Token,
	}, nil
}
