// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleCategory

import (
	"context"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/models"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateArticleCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleCategoryLogic {
	return &UpdateArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateArticleCategoryLogic) UpdateArticleCategory(req *types.UpdateArticleCategoryReq) (resp *types.Empty, err error) {
	if req.Id <= 0 {
		return nil, errorx.ErrorInvalidParam
	}
	uid := middleware.GetUidFromContext(l.ctx)
	category := models.ArticleCategory{
		Model: gorm.Model{
			ID: uint(req.Id),
		},
		Uid:         uint(uid),
		Name:        req.Name,
		Description: req.Description,
	}
	if err := l.svcCtx.ArticleCategoryModel.Update(l.ctx, category); err != nil {
		l.Errorf("更新分类失败: %v", err)
		return nil, err
	}

	return &types.Empty{}, nil
}
