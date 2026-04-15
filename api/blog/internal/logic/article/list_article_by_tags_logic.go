// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticleByTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListArticleByTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticleByTagsLogic {
	return &ListArticleByTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListArticleByTagsLogic) ListArticleByTags(req *types.ListArticleByTagsReq) (resp *types.ListArticleByTagsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
