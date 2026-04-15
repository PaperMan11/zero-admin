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

type DeregistrationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeregistrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeregistrationLogic {
	return &DeregistrationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeregistrationLogic) Deregistration(req *types.Empty) (resp *types.Empty, err error) {
	tokenUUID := middleware.GetTokenIDFromContext(l.ctx)
	userId := middleware.GetUidFromContext(l.ctx)
	if tokenUUID != "" && userId != 0 {
		accessKey := utils.GetAccessTokenKey(userId, tokenUUID)
		refreshKey := utils.GetRefreshTokenKey(userId, tokenUUID)
		l.svcCtx.LocalCache.Del(accessKey)
		l.svcCtx.LocalCache.Del(refreshKey)
		l.svcCtx.Redis.DelCtx(l.ctx, accessKey, refreshKey)
	}

	user, err := l.svcCtx.UserModel.GetByID(l.ctx, uint(userId))
	if err != nil {
		l.Errorf("检查用户状态失败: %v, 用户ID: %d", err, userId)
		return nil, err
	}

	// 注销用户
	err = l.svcCtx.UserModel.Delete(l.ctx, user.ID)
	if err != nil {
		l.Errorf("注销用户失败: %v, 用户ID: %d", err, userId)
		return nil, err
	}

	return &types.Empty{}, nil
}
