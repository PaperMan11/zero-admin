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

func (l *RegisterLogic) Register(req *types.RegisterRequest, ip, os, browser string) (resp *types.RegisterResponse, err error) {
	res, err := l.svcCtx.AuthService.Register(l.ctx, &authservice.RegisterRequest{
		Username:  req.Username,
		Password:  req.Password,
		IpAddress: ip,
		Email:     req.Email,
		Mobile:    req.Mobile,
		RealName:  req.RealName,
		Gender:    req.Gender,
		Avatar:    req.Avatar,
		Os:        os,
		Browser:   browser,
	})
	if err != nil {
		logc.Errorf(l.ctx, "用户注册：%+v,异常:%s", req, err.Error())
		return nil, err
	}

	// 添加token过期管理
	SetTokenCache(l.ctx, l.svcCtx, res.Id, res.TokenUuid, res.RefreshToken)

	return &types.RegisterResponse{
		Id:           res.Id,
		RefreshToken: res.RefreshToken,
		Token:        res.Token,
		Username:     res.Username,
	}, nil
}
