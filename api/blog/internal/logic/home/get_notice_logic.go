// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package home

import (
	"context"
	"zero-admin/api/blog/internal/models"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeLogic {
	return &GetNoticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNoticeLogic) GetNotice(req *types.Empty) (resp *types.Notice, err error) {
	n, err := l.svcCtx.Barrier.Do("GetNotice", func() (any, error) {
		return l.svcCtx.NoticeModel.GetLatest(l.ctx)
	})
	if err != nil {
		l.Errorf("获取最新公告失败: %v", err)
		return nil, errorx.ErrorInternal
	}
	notice, ok := n.(*models.Notice)
	if !ok {
		l.Errorf("获取最新公告失败: %v", err)
		return nil, errorx.ErrorInternal
	}
	return &types.Notice{
		Content:   notice.Content,
		CreatedAt: notice.CreatedAt.Unix(),
		Id:        int64(notice.ID),
		UpdatedAt: notice.UpdatedAt.Unix(),
	}, nil
}
