// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleComments

import (
	"context"
	"math"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplyCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReplyCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplyCommentsLogic {
	return &GetReplyCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReplyCommentsLogic) GetReplyComments(req *types.GetReplyCommentsReq) (resp *types.GetReplyCommentsResp, err error) {
	comment, err := l.svcCtx.ArticleCommentModel.GetByID(l.ctx, uint(req.CommentId))
	if err != nil {
		l.Errorf("获取文章评论总数失败: %v", err)
		return nil, errorx.ErrorInternal
	}

	page := req.Page
	pageSize := req.PageSize
	l.Debugf("获取文章评论总数: %d, page: %d, pageSize: %d", comment.ReplyCount, page, pageSize)
	commentReplies, err := l.svcCtx.ArticleCommentModel.ListReplyComments(l.ctx, req.ArticleId, req.CommentId, int(page), int(pageSize))
	if err != nil {
		l.Errorf("获取评论回复失败: %v", err)
		return nil, errorx.ErrorInternal
	}

	respComments := make([]types.Reply, 0, len(commentReplies))
	for _, reply := range commentReplies {
		respComments = append(respComments, types.Reply{
			Id:              int64(reply.ID),
			UserId:          reply.UserId,
			UserName:        reply.UserName,
			UserAvatar:      reply.UserAvatar,
			Content:         reply.Content,
			ReplyUserId:     reply.ReplyTo,
			ReplyUserName:   reply.ReplyToName,
			ReplyUserAvatar: reply.ReplyToAvatar,
			CreatedAt:       reply.CreatedAt.Unix(),
		})
	}
	return &types.GetReplyCommentsResp{
		PageResponse: types.PageResponse{
			Total:     comment.ReplyCount,
			Page:      page,
			PageSize:  pageSize,
			TotalPage: int64(math.Ceil(float64(comment.ReplyCount) / float64(pageSize))),
		},
		Replies: respComments,
	}, nil
}
