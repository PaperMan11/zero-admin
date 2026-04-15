// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"
	"strings"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleLogic) GetArticle(req *types.IdValue) (resp *types.Article, err error) {
	article, err := l.svcCtx.ArticleModel.GetDetailByID(l.ctx, uint(req.Id))
	if err != nil {
		l.Errorf("查询文章详情失败: %v, 文章ID: %d", err, req.Id)
		return nil, err
	}

	return &types.Article{
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
	}, nil
}
