// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"
	"strings"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticleByCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListArticleByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticleByCategoryLogic {
	return &ListArticleByCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListArticleByCategoryLogic) ListArticleByCategory(req *types.ListArticleByCategoryReq) (resp *types.ListArticleByCategoryResp, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	articles, err := l.svcCtx.ArticleModel.GetByUidAndCategory(l.ctx, uint(uid), req.CategoryId, int(req.Page), int(req.PageSize))
	if err != nil {
		l.Errorf("获取分类文章失败: %d, %v", req.CategoryId, err)
		return nil, err
	}

	total, err := l.svcCtx.ArticleModel.GetTotalByCategoryId(l.ctx, req.CategoryId)
	if err != nil {
		l.Errorf("获取分类文章总数量失败: %d, %v", req.CategoryId, err)
		return nil, err
	}

	respArticles := make([]types.Article, 0, len(articles))
	for _, article := range articles {
		respArticles = append(respArticles, types.Article{
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
	return &types.ListArticleByCategoryResp{
		PageResponse: types.PageResponse{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: total/req.PageSize + 1,
		},
		Articles: respArticles,
	}, nil
}
