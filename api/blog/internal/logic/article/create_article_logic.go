// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"
	"strings"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/models"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/convert"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleLogic {
	return &CreateArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateArticleLogic) CreateArticle(req *types.CreateArticleReq) (resp *types.CreateArticleResp, err error) {
	authorId := middleware.GetUidFromContext(l.ctx)
	// 1.判断用户是否存在
	user, err := l.svcCtx.UserModel.GetByID(l.ctx, uint(authorId))
	if err != nil {
		l.Errorf("检查用户状态失败: %v, 用户ID: %d", err, authorId)
		return nil, err
	}

	// 2.判断分类是否存在
	_, err = l.svcCtx.ArticleCategoryModel.GetByID(l.ctx, uint(req.CategoryId))
	if err != nil {
		l.Errorf("检查分类是否存在失败: %v, 分类ID: %d", err, req.CategoryId)
		return nil, err
	}

	article := models.Article{
		AuthorId:   convert.ToInt64(user.ID),
		CategoryId: req.CategoryId,
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Tags:       strings.Join(req.Tags, ","),
		Cover:      req.Cover,
	}
	articleId, err := l.svcCtx.ArticleModel.Create(l.ctx, article)
	if err != nil {
		l.Errorf("创建文章失败: %v", err)
		return nil, err
	}

	// 更新分类文章数量
	_ = l.svcCtx.ArticleCategoryModel.IncArticleCount(l.ctx, uint(req.CategoryId))

	return &types.CreateArticleResp{
		Id: int64(articleId),
	}, nil
}
