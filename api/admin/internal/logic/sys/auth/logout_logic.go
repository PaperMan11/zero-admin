// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/pkg/convert"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() (resp *types.Empty, err error) {
	// 添加token过期管理
	uid := convert.ToInt64(l.ctx.Value("uid"))
	l.svcCtx.Redis.DelCtx(l.ctx, utils.GetAccessTokenKey(uid), utils.GetRefreshTokenKey(uid))
	return
}
