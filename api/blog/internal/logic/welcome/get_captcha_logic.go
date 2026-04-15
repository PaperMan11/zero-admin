// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package welcome

import (
	"context"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

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

func (l *GetCaptchaLogic) GetCaptcha(req *types.Empty) (resp *types.CaptchaResp, err error) {
	id, s, _, err := l.svcCtx.Captcha.Generate()
	if err != nil {
		l.Errorf("生成验证码失败: %v", err)
		return nil, err
	}
	return &types.CaptchaResp{
		CaptchaId:    id,
		CaptchaImage: s,
	}, nil
}
