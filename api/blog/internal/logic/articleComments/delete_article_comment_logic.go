// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleComments

import (
	"context"

	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteArticleCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleCommentLogic {
	return &DeleteArticleCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteArticleCommentLogic) DeleteArticleComment(req *types.DeleteArticleCommentReq) (resp *types.Empty, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	comment, err := l.svcCtx.ArticleCommentModel.GetByID(l.ctx, uint(req.CommentId))
	if err != nil {
		l.Errorf("获取文章评论失败: %v, 文章评论ID: %d", err, req.CommentId)
		return nil, err
	}

	article, err := l.svcCtx.ArticleModel.GetByID(l.ctx, uint(req.ArticleId))
	if err != nil {
		l.Errorf("获取文章失败: %v, 文章ID: %d", err, req.ArticleId)
		return nil, err
	}

	if uid != comment.UserId && uid != article.AuthorId {
		l.Errorf("您没有权限删除该文章评论, 文章评论ID: %d, 文章ID: %d", req.CommentId, req.ArticleId)
		return nil, errorx.ErrorArticleCommentNoPermission
	}

	if err := l.svcCtx.ArticleCommentModel.Delete(l.ctx, uint(req.CommentId)); err != nil {
		l.Errorf("删除文章评论失败: %v, 文章评论ID: %d", err, req.CommentId)
		return nil, errorx.ErrorInternal
	}
	return &types.Empty{}, nil
}
