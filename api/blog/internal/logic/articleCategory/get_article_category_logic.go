// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleCategory

import (
	"context"
	"strings"
	"zero-admin/api/blog/internal/models"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleCategoryLogic {
	return &GetArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleCategoryLogic) GetArticleCategory(req *types.IdValue) (resp *types.ArticleCategory, err error) {
	uid := middleware.GetUidFromContext(l.ctx)

	var category *models.ArticleCategory
	if req.Id == 0 {
		total, _ := l.svcCtx.ArticleModel.CountArticleByCategoryId(l.ctx, uint(uid), 0)
		category = &models.ArticleCategory{
			Uid:          uint(uid),
			Name:         "默认",
			Description:  "默认分类",
			ArticleCount: total,
		}
	} else {
		category, err = l.svcCtx.ArticleCategoryModel.GetByID(l.ctx, uint(req.Id))
		if err != nil {
			l.Errorf("获取文章分类失败: %v", err)
			return nil, err
		}
	}

	// 获取最新5篇文章
	articles, err := l.svcCtx.ArticleModel.GetByUidAndCategory(l.ctx, uint(uid), int64(category.ID), 1, 5)
	if err != nil {
		l.Errorf("获取文章失败: %v", err)
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
			CreatedTime: article.CreatedAt.Unix(),
			UpdatedTime: article.UpdatedAt.Unix(),
		})
	}
	return &types.ArticleCategory{
		Id:           int64(category.ID),
		Name:         category.Name,
		Description:  category.Description,
		Articles:     respArticles,
		ArticleCount: category.ArticleCount,
	}, nil
}
