// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaLogic {
	return &GetCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCaptchaLogic) GetCaptcha() (resp *types.CaptchaResponse, err error) {
	id, s, _, err := l.svcCtx.Captcha.Generate()
	if err != nil {
		logc.Errorf(l.ctx, "生成验证码失败: %v", err)
		return nil, err
	}
	return &types.CaptchaResponse{
		CaptchaId:    id,
		CaptchaImage: s,
	}, nil
}
