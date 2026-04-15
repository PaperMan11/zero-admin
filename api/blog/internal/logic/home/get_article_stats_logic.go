// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package home

import (
	"context"
	"zero-admin/api/blog/internal/middleware"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleStatsLogic {
	return &GetArticleStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleStatsLogic) GetArticleStats(req *types.Empty) (resp *types.ArticleStats, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	var articleNum, categoryNum, commentNum int64
	articleNum, err = l.svcCtx.ArticleModel.CountArticleByUid(l.ctx, uint(uid))
	if err != nil {
		l.Errorf("获取用户总文章数量失败 uid: %d, err: %v", uid, err)
		return nil, errorx.ErrorInternal
	}
	categoryNum, err = l.svcCtx.ArticleCategoryModel.CountArticleCategoryByUid(l.ctx, uint(uid))
	if err != nil {
		l.Errorf("获取用户总分类数量失败 uid: %d, err: %v", uid, err)
		return nil, errorx.ErrorInternal
	}
	commentNum, err = l.svcCtx.ArticleCommentModel.CountCommentByUserId(l.ctx, uint(uid))
	if err != nil {
		l.Errorf("获取用户总评论数量失败 uid: %d, err: %v", uid, err)
		return nil, errorx.ErrorInternal
	}
	return &types.ArticleStats{
		ArticleNum:  articleNum,
		CategoryNum: categoryNum,
		CommentNum:  commentNum,
	}, nil
}
