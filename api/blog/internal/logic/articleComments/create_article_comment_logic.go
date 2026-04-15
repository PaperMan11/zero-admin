// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleComments

import (
	"context"
	"strings"

	"zero-admin/api/blog/internal/models"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateArticleCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateArticleCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleCommentLogic {
	return &CreateArticleCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateArticleCommentLogic) CreateArticleComment(req *types.CreateArticleCommentReq) (resp *types.CreateArticleCommentResp, err error) {
	comment := models.ArticleComment{
		ArticleId: req.ArticleId,
		ParentId:  req.ParentCommentId, // 回复评论的父评论id
		ReplyId:   req.ReplyCommentId,  // 回复的评论id
		UserId:    req.UserId,
		ReplyTo:   req.ReplyTo, // 回复用户ID
		Content:   req.Content,
		Images:    strings.Join(req.Imgs, ","), // 图片
	}
	id, err := l.svcCtx.ArticleCommentModel.Create(l.ctx, comment)
	if err != nil {
		l.Errorf("创建评论失败: %v", err)
		return nil, errorx.ErrorInternal
	}
	if req.ParentCommentId > 0 {
		l.svcCtx.ArticleCommentModel.IncCommentReplyCount(l.ctx, uint(req.ParentCommentId), 1)
		if req.ParentCommentId != req.ReplyCommentId {
			l.svcCtx.ArticleCommentModel.IncCommentReplyCount(l.ctx, uint(req.ReplyCommentId), 1)
		}
	}
	return &types.CreateArticleCommentResp{
		Id: int64(id),
	}, nil
}
