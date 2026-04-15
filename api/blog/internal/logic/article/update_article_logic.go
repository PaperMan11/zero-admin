// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"
	"strings"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleLogic {
	return &UpdateArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateArticleLogic) UpdateArticle(req *types.UpdateArticleReq) (resp *types.Article, err error) {
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.CategoryId != 0 {
		updates["category_id"] = req.CategoryId
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Tags != nil {
		tags := strings.Join(req.Tags, ",")
		updates["tags"] = tags
	}
	if err := l.svcCtx.ArticleModel.Updates(l.ctx, uint(req.Id), updates); err != nil {
		l.Errorf("更新文章信息失败: %v, 文章ID: %d", err, req.Id)
		return nil, errorx.ErrorInternal
	}

	article, err := l.svcCtx.ArticleModel.GetDetailByID(l.ctx, uint(req.Id))
	if err != nil {
		l.Errorf("检查文章状态失败: %v, 文章ID: %d", err, req.Id)
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
