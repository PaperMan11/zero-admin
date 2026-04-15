// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleCategory

import (
	"context"
	"strings"
	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticleCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticleCategoryLogic {
	return &ListArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListArticleCategoryLogic) ListArticleCategory(req *types.ListArticleCategoryReq) (resp *types.ListArticleCategoryResp, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	res, err := l.svcCtx.ArticleCategoryModel.PaginationByUid(l.ctx, uint(uid), int(req.Page), int(req.PageSize), req.FullQuery)
	if err != nil {
		l.Errorf("获取用户文章分类失败 uid: %d, err: %s", uid, err)
		return nil, errorx.ErrorInternal
	}

	categories := make([]types.ArticleCategory, 0, len(res)+1)
	// 默认分类
	articles, _ := l.svcCtx.ArticleModel.GetByUidAndCategory(l.ctx, uint(uid), 0, 1, 5)
	total, _ := l.svcCtx.ArticleModel.GetTotalByCategoryId(l.ctx, 0)
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
			CreatedTime: article.CreatedAt.Unix(),
			UpdatedTime: article.UpdatedAt.Unix(),
		})
	}
	categories = append(categories, types.ArticleCategory{
		Id:           0,
		Name:         "默认",
		Description:  "默认",
		Articles:     respArticles,
		ArticleCount: total,
	})
	for _, v := range res {
		// 获取最新5篇文章
		articles, err := l.svcCtx.ArticleModel.GetByUidAndCategory(l.ctx, v.Uid, int64(v.ID), 1, 5)
		if err != nil {
			l.Errorf("获取文章失败 uid: %d, categoryId:  err: %v", uid, v.ID, err)
			continue
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
				CreatedTime: article.CreatedAt.Unix(),
				UpdatedTime: article.UpdatedAt.Unix(),
			})
		}
		categories = append(categories, types.ArticleCategory{
			Id:           int64(v.ID),
			Name:         v.Name,
			Description:  v.Description,
			Articles:     respArticles,
			ArticleCount: v.ArticleCount,
		})
	}
	categoryTotal, _ := l.svcCtx.ArticleCategoryModel.CountArticleCategoryByUid(l.ctx, uint(uid))
	return &types.ListArticleCategoryResp{
		PageResponse: types.PageResponse{
			Total:     categoryTotal,
			Page:      req.Page,
			PageSize:  req.PageSize,
			TotalPage: total/req.PageSize + 1,
		},
		Categories: categories,
	}, nil
}
