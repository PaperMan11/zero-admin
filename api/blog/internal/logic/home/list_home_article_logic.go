// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package home

import (
	"context"
	"strings"
	"zero-admin/api/blog/internal/models"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListHomeArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListHomeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListHomeArticleLogic {
	return &ListHomeArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListHomeArticleLogic) ListHomeArticle(req *types.HomeArticleReq) (resp *types.HomeArticleResp, err error) {
	var articles []models.ArticleDetail
	var respArticles []types.Article
	switch req.Tab {
	case "recent":
		articles, err = l.svcCtx.ArticleModel.ListRecentArticle(l.ctx, int(req.Page), int(req.PageSize))
	case "hot":
		articles, err = l.svcCtx.ArticleModel.ListHotArticle(l.ctx, int(req.Page), int(req.PageSize))
	case "vote":
		articles, err = l.svcCtx.ArticleModel.ListVoteArticle(l.ctx, int(req.Page), int(req.PageSize))
	default:
		return nil, errorx.ErrorInvalidParam
	}
	if err != nil {
		l.Errorf("获取文章列表失败 req: %+v, err: %v", req, err)
		return nil, errorx.ErrorInternal
	}
	for _, article := range articles {
		respArticles = append(respArticles, types.Article{
			Id:         int64(article.ID),
			AuthorId:   article.AuthorId,
			AuthorName: article.AuthorName,
			Title:      article.Title,
			Summary:    article.Summary,
			//Content:     article.Content,
			Views:       article.ViewCount,
			Likes:       article.LikeCount,
			Comments:    article.CommentCount,
			Cover:       article.Cover,
			Tags:        strings.Split(article.Tags, ","),
			CreatedTime: article.CreatedAt.Unix(),
			UpdatedTime: article.UpdatedAt.Unix(),
		})
	}

	total, err := l.svcCtx.ArticleModel.GetTotal(l.ctx)
	if err != nil {
		l.Errorf("获取文章总数失败 req: %+v, err: %v", req, err)
		return nil, errorx.ErrorInternal
	}

	return &types.HomeArticleResp{
		Articles: respArticles,
		PageResponse: types.PageResponse{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	}, nil
}
