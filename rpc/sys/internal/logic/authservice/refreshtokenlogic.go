package authservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/pkg/convert"
	jwtUtil "zero-admin/pkg/jwt"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 刷新令牌
func (l *RefreshTokenLogic) RefreshToken(in *sysclient.RefreshTokenRequest) (*sysclient.RefreshTokenResponse, error) {
	issuer := l.svcCtx.Config.Name
	accessSecret := l.svcCtx.Config.Auth.AccessSecret
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	refreshSecret := l.svcCtx.Config.Auth.RefreshSecret
	refreshExpire := l.svcCtx.Config.Auth.RefreshExpire

	userIDStr, err := jwtUtil.ParseRefreshToken(issuer, in.RefreshToken, refreshSecret)
	if err != nil {
		logc.Debug(l.ctx, "refresh token error: %v", err)
		return nil, status.Error(codes.Unauthenticated, "请重新登录")
	}
	userID := convert.ToInt64(userIDStr)

	user, err := l.svcCtx.DB.GetUserByID(l.ctx, userID)
	// 判断用户是否存在
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		logc.Errorf(l.ctx, "用户不存在 [userID:%d]", userID)
		return nil, errors.New("用户不存在")
	case err != nil:
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "查询用户信息异常")
	}

	// 用户角色信息
	roleCodes, _ := l.svcCtx.DB.GetUserRoleCodes(l.ctx, user.ID)
	accessToken, refreshToken, err := GenerateToken(user.ID, roleCodes, issuer, accessSecret, accessExpire, refreshSecret, refreshExpire)
	if err != nil {
		logc.Errorf(l.ctx, "生成token异常, 用户id：%+v, 错误：%s", user.ID, err.Error())
		return nil, status.Error(codes.Internal, "生成token异常")
	}

	return &sysclient.RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
