// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteArticleLogic) DeleteArticle(req *types.DeleteArticleReq) (resp *types.Empty, err error) {
	// 1.判断文章是否存在
	article, err := l.svcCtx.ArticleModel.GetByID(l.ctx, uint(req.Id))
	if err != nil {
		l.Errorf("检查文章状态失败: %v, 文章ID: %d", err, req.Id)
		return nil, err
	}

	// 2.删除文章
	err = l.svcCtx.ArticleModel.Delete(l.ctx, uint(req.Id))
	if err != nil {
		l.Errorf("删除文章失败: %v, 文章ID: %d", err, req.Id)
		return nil, err
	}

	// 3.减少文章分类文章数量
	if article.CategoryId > 0 {
		_ = l.svcCtx.ArticleCategoryModel.DecArticleCount(l.ctx, uint(article.CategoryId))
	}

	return &types.Empty{}, nil
}
