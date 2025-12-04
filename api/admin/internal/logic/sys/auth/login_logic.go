// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
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
	return &types.LoginResponse{
		AccessToken:  res.Token,
		Id:           res.Id,
		RefreshToken: res.RefreshToken,
		Username:     res.Username,
	}, nil
}
