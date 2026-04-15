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

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchResp, err error) {
	articles, err := l.svcCtx.ArticleModel.Search(l.ctx, req.Keyword, int(req.Page), int(req.PageSize))
	if err != nil {
		l.Errorf("搜索失败 keywords: %s, err: %s", req.Keyword, err)
		return nil, errorx.ErrorInternal
	}
	results := make([]types.Article, 0, len(articles))
	for _, article := range articles {
		results = append(results, types.Article{
			Id:          int64(article.ID),
			AuthorId:    article.AuthorId,
			AuthorName:  article.AuthorName,
			Title:       article.Title,
			Summary:     article.Summary,
			Content:     article.Content,
			Views:       article.ViewCount,
			Likes:       article.LikeCount,
			Comments:    article.CommentCount,
			Tags:        strings.Split(article.Tags, ","),
			Cover:       article.Cover,
			CreatedTime: article.CreatedAt.Unix(),
			UpdatedTime: article.UpdatedAt.Unix(),
		})
	}
	total, _ := l.svcCtx.ArticleModel.CountArticleByKeyword(l.ctx, req.Keyword)

	return &types.SearchResp{
		PageResponse: types.PageResponse{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: total/req.PageSize + 1,
		},
		Results: results,
	}, nil
}
