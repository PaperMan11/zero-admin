// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleComments

import (
	"context"
	"math"
	"zero-admin/api/blog/internal/middleware"
	"zero-admin/pkg/errorx"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnreadCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUnreadCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadCommentsLogic {
	return &GetUnreadCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnreadCommentsLogic) GetUnreadComments(req *types.GetUnreadCommentsReq) (resp *types.GetUnreadCommentsResp, err error) {
	uid := middleware.GetUidFromContext(l.ctx)
	total, _ := l.svcCtx.ArticleCommentModel.CountUnreadComments(l.ctx, uid)
	comments, err := l.svcCtx.ArticleCommentModel.PaginationUnreadComments(l.ctx, uid, int(req.Page), int(req.PageSize))
	if err != nil {
		l.Errorf("获取未读评论失败 uid:%d, err:%v", uid, err)
		return nil, errorx.ErrorInternal
	}

	respComments := make([]types.Reply, 0, len(comments))
	for _, reply := range comments {
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
	return &types.GetUnreadCommentsResp{
		PageResponse: types.PageResponse{
			Total:     total,
			Page:      req.Page,
			PageSize:  req.PageSize,
			TotalPage: int64(math.Ceil(float64(total) / float64(req.PageSize))),
		},
		Replies: respComments,
	}, nil
}
