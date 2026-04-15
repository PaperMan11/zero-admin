// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package welcome

import (
	"context"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoffLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoffLogic {
	return &LogoffLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoffLogic) Logoff(req *types.Empty) (resp *types.Empty, err error) {
	tokenUUID := middleware.GetTokenIDFromContext(l.ctx)
	userId := middleware.GetUidFromContext(l.ctx)

	if tokenUUID != "" && userId != 0 {
		accessKey := utils.GetAccessTokenKey(userId, tokenUUID)
		refreshKey := utils.GetRefreshTokenKey(userId, tokenUUID)
		l.svcCtx.LocalCache.Del(accessKey)
		l.svcCtx.LocalCache.Del(refreshKey)
		l.svcCtx.Redis.DelCtx(l.ctx, accessKey, refreshKey)
	}
	return &types.Empty{}, nil
}
