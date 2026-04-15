// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"strings"
	"zero-admin/api/blog/internal/middleware"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserArticleTimelineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserArticleTimelineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserArticleTimelineLogic {
	return &GetUserArticleTimelineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserArticleTimelineLogic) GetUserArticleTimeline(req *types.Empty) (resp []types.Article, err error) {
	uid := middleware.GetUidFromContext(l.ctx)

	articles, err := l.svcCtx.ArticleModel.GetTimeline(l.ctx, uid)
	if err != nil {
		l.Errorf("获取用户时间线失败: %v", err)
		return nil, errorx.ErrorInternal
	}

	resp = make([]types.Article, 0, len(articles))
	for _, article := range articles {
		resp = append(resp, types.Article{
			Id:          int64(article.ID),
			AuthorId:    article.AuthorId,
			Title:       article.Title,
			Summary:     article.Summary,
			Content:     article.Content,
			Views:       article.ViewCount,
			Likes:       article.LikeCount,
			Cover:       article.Cover,
			Comments:    article.CommentCount,
			Tags:        strings.Split(article.Tags, ","),
			CreatedTime: article.CreatedAt.Unix(),
			UpdatedTime: article.UpdatedAt.Unix(),
		})
	}
	return
}
