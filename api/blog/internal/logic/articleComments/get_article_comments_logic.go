// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package articleComments

import (
	"context"
	"math"
	"strings"
	"zero-admin/api/blog/internal/models"

	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"
	"zero-admin/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleCommentsLogic {
	return &GetArticleCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleCommentsLogic) GetArticleComments(req *types.GetArticleCommentsReq) (resp *types.GetArtcileCommentsResp, err error) {
	total, err := l.svcCtx.ArticleCommentModel.GetTotalByArticleId(l.ctx, req.ArticleId)
	if err != nil {
		l.Errorf("获取文章评论总数失败: %v", err)
		return nil, errorx.ErrorInternal
	}

	page := req.Page
	pageSize := req.PageSize
	replyPage := 0
	commentId := int64(0)
	if req.TargetCommentId > 0 {
		comment, _ := l.svcCtx.ArticleCommentModel.GetByID(l.ctx, uint(req.TargetCommentId))
		if comment.ID > 0 {
			commentId = comment.ParentId
			if comment.ParentId == 0 {
				// 获取指定一级评论id是第几条评论
				index, _ := l.svcCtx.ArticleCommentModel.GetCommentIndex(l.ctx, uint(req.TargetCommentId))
				// 根据评论索引计算页码
				page = int64(int(index-1)/int(pageSize) + 1)
			} else {
				// 获取指定一级评论id是第几条评论
				index, _ := l.svcCtx.ArticleCommentModel.GetCommentIndex(l.ctx, uint(comment.ParentId))
				page = int64(int(index-1)/int(pageSize) + 1)
				// 获取指定二级评论id是第几条评论
				index2, _ := l.svcCtx.ArticleCommentModel.GetReplyCommentIndex(l.ctx, uint(comment.ParentId), comment.ID)
				replyPage = int(index2-1)/int(pageSize) + 1
			}
		}
	}
	comments, err := l.svcCtx.ArticleCommentModel.ListParentComment(l.ctx, req.ArticleId, req.Tab, int(page), int(pageSize))
	if err != nil {
		l.Errorf("获取文章评论失败: %v", err)
		return nil, errorx.ErrorInternal
	}
	var replies []models.ArticleCommentReply
	if commentId > 0 && replyPage > 0 {
		replies, _ = l.svcCtx.ArticleCommentModel.ListReplyComments(l.ctx, req.ArticleId, commentId, 1, replyPage*int(pageSize))
	}
	//l.Infof("获取文章评论成功: 一级评论列表: %v, 二级评论列表: %v, %d, replyPageSize:%d", comments, replies, replyPageSize, replyPageSize)

	respReplies := make([]types.Reply, 0, len(replies))
	for _, reply := range replies {
		respReplies = append(respReplies, types.Reply{
			Id:              int64(reply.ID),
			UserId:          reply.UserId,
			UserName:        reply.UserName,
			UserAvatar:      reply.UserAvatar,
			Content:         reply.Content,
			ReplyUserId:     reply.ReplyTo,
			ReplyUserName:   reply.ReplyToName,
			ReplyUserAvatar: reply.ReplyToAvatar,
			CreatedAt:       reply.CreatedAt.Unix(),
			Read:            int(reply.Read),
		})
	}

	respComments := make([]types.Comment, 0, len(comments))
	for _, comment := range comments {
		respComment := types.Comment{
			Id:         int64(comment.ID),
			UserId:     comment.UserId,
			UserName:   comment.UserName,
			UserAvatar: comment.UserAvatar,
			Content:    comment.Content,
			Images:     strings.Split(comment.Images, ","),
			Replies:    nil,
			ReplyCount: comment.ReplyCount,
			CreatedAt:  comment.CreatedAt.Unix(),
		}
		if comment.ID == uint(commentId) {
			respComment.Replies = respReplies
		}
		respComments = append(respComments, respComment)
	}
	return &types.GetArtcileCommentsResp{
		PageResponse: types.PageResponse{
			Total:     total,
			Page:      page,
			PageSize:  pageSize,
			TotalPage: int64(math.Ceil(float64(total) / float64(pageSize))),
		},
		Comments: respComments,
	}, nil
}
