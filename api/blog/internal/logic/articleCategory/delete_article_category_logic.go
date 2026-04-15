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

type DeleteArticleCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleCategoryLogic {
	return &DeleteArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteArticleCategoryLogic) DeleteArticleCategory(req *types.DeleteArticleCategoryReq) (resp *types.Empty, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	err = l.svcCtx.ArticleCategoryModel.Delete(l.ctx, uint(req.Id), uint(uid))
	if err != nil {
		l.Errorf("删除文章分类失败: %v", err)
		return nil, err
	}
	return &types.Empty{}, nil
}
