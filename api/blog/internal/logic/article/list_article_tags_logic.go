// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticleTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListArticleTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticleTagsLogic {
	return &ListArticleTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListArticleTagsLogic) ListArticleTags(req *types.ListArticleTagsReq) (resp *types.ListArticleTagsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
