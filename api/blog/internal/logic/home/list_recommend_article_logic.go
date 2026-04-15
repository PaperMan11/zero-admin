// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package home

import (
	"context"
	"strings"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRecommendArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRecommendArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRecommendArticleLogic {
	return &ListRecommendArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRecommendArticleLogic) ListRecommendArticle(req *types.Empty) (resp []types.Article, err error) {
	articles, err := l.svcCtx.ArticleModel.GetRandomArticle(l.ctx)
	if err != nil {
		l.Errorf("获取推荐文章失败: %v", err)
		return nil, errorx.ErrorInternal
	}
	resp = make([]types.Article, 0, len(articles))
	for _, article := range articles {
		resp = append(resp, types.Article{
			Id:         int64(article.ID),
			AuthorId:   article.AuthorId,
			AuthorName: article.AuthorName,
			Title:      article.Title,
			Summary:    article.Summary,
			//Content:     article.Content,
			Views:       article.ViewCount,
			Likes:       article.LikeCount,
			Cover:       article.Cover,
			Comments:    article.CommentCount,
			Tags:        strings.Split(article.Tags, ","),
			CreatedTime: article.CreatedAt.Unix(),
			UpdatedTime: article.UpdatedAt.Unix(),
		})
	}
	return resp, nil
}
