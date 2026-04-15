// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleCategory

import (
	"context"
	"zero-admin/api/blog/internal/middleware"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllArticleCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllArticleCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllArticleCategoriesLogic {
	return &GetAllArticleCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllArticleCategoriesLogic) GetAllArticleCategories(req *types.Empty) (resp []types.ArticleCategory, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	categories, err := l.svcCtx.ArticleCategoryModel.ListAllByUid(l.ctx, uint(uid))
	if err != nil {
		l.Errorf("获取用户文章分类失败 uid: %d, err: %s", uid, err)
		return
	}
	resp = make([]types.ArticleCategory, 0, len(categories))
	for _, category := range categories {
		resp = append(resp, types.ArticleCategory{
			Id:          int64(category.ID),
			Name:        category.Name,
			Description: category.Description,
		})
	}
	return resp, nil
}
