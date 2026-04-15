// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleCategory

import (
	"context"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/models"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateArticleCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleCategoryLogic {
	return &CreateArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateArticleCategoryLogic) CreateArticleCategory(req *types.CreateArticleCategoryReq) (resp *types.CreateArticleCategoryResp, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	exist, err := l.svcCtx.ArticleCategoryModel.CheckUserCategoryNameExist(l.ctx, uint(uid), req.Name)
	if err != nil {
		l.Errorf("检查分类名是否存在失败: %v", err)
		return nil, err
	}
	if exist {
		return nil, errorx.ErrorCategoryNameExist
	}

	id, err := l.svcCtx.ArticleCategoryModel.Create(l.ctx, models.ArticleCategory{
		Uid:         uint(uid),
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		l.Errorf("创建分类失败: %v", err)
		return nil, err
	}

	return &types.CreateArticleCategoryResp{
		Id: int64(id),
	}, nil
}
