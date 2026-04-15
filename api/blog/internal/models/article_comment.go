package models

import (
	"context"
	"errors"
	"fmt"
	"time"
	"zero-admin/api/blog/internal/models/generated"
	"zero-admin/pkg/errorx"

	"gorm.io/gorm"
)

type ArticleComment struct {
	gorm.Model
	ArticleId  int64  `json:"articleId" gorm:"column:article_id"`
	ParentId   int64  `json:"parentId" gorm:"column:parent_id"` // 父评论ID，0表示一级评论
	ReplyId    int64  `json:"replyId" gorm:"column:reply_id"`   // 回复评论ID
	UserId     int64  `json:"userId" gorm:"column:user_id"`
	ReplyTo    int64  `json:"replyTo" gorm:"column:reply_to"` // 回复用户ID
	Content    string `json:"content" gorm:"column:content"`
	Images     string `json:"images" gorm:"column:images"`          // 图片
	ReplyCount int64  `json:"replyCount" gorm:"column:reply_count"` // 回复数
	Read       int8   `json:"read" gorm:"column:read"`              // 0:未读 1:已读
}

func (ArticleComment) TableName() string {
	return "blog_article_comment"
}

type ArticleCommentDetail struct {
	ArticleComment
	UserName   string `json:"userName" gorm:"column:user_name"`
	UserAvatar string `json:"userAvatar" gorm:"column:user_avatar"`
}

func (ArticleCommentDetail) TableName() string {
	return "blog_article_comment"
}

type ArticleCommentReply struct {
	ArticleComment
	UserName      string `json:"userName" gorm:"column:user_name"`
	UserAvatar    string `json:"userAvatar" gorm:"column:user_avatar"`
	ReplyToName   string `json:"replyToName" gorm:"column:reply_to_name"`
	ReplyToAvatar string `json:"replyToAvatar" gorm:"column:reply_to_avatar"`
}

type ArticleCommentModel struct {
	db *gorm.DB
}

func NewArticleCommentModel(db *gorm.DB) *ArticleCommentModel {
	return &ArticleCommentModel{db}
}

func (m *ArticleCommentModel) Create(ctx context.Context, comment ArticleComment) (id uint, err error) {
	err = gorm.G[ArticleComment](m.db).Create(ctx, &comment)
	return comment.ID, err
}

func (m *ArticleCommentModel) Delete(ctx context.Context, id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[ArticleComment](tx).Where(generated.ArticleComment.ID.Eq(id)).
			Set(generated.ArticleComment.DeletedAt.Set(gorm.DeletedAt{Time: time.Now(), Valid: true})).Update(ctx)
		if err != nil {
			return err
		}

		_, err = gorm.G[ArticleComment](tx).Where(generated.ArticleComment.ParentId.Eq(int64(id))).
			Set(generated.ArticleComment.DeletedAt.Set(gorm.DeletedAt{Time: time.Now(), Valid: true})).Update(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (m *ArticleCommentModel) GetByID(ctx context.Context, id uint) (comment *ArticleComment, err error) {
	var c ArticleComment
	c, err = gorm.G[ArticleComment](m.db).Where(generated.ArticleComment.ID.Eq(id)).
		Where(generated.Article.DeletedAt.IsNull()).
		Where(generated.ArticleComment.DeletedAt.IsNull()).
		First(ctx)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errorx.ErrorArticleCommentNotFound
	case err != nil:
		return nil, errorx.ErrorInternal
	}
	return &c, err
}

func (m *ArticleCommentModel) ListParentComment(ctx context.Context, articleId int64, tab string, page, pageSize int) (comments []ArticleCommentDetail, err error) {
	sql := m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS c", ArticleComment{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON c.user_id = u.id", User{}.TableName())).
		Where("c.parent_id = 0").
		Where("c.article_id = ?", articleId).
		Where("c.deleted_at IS NULL")

	if tab == "hot" {
		sql.Order("c.reply_count DESC, c.id DESC")
	} else {
		sql.Order("c.id DESC")
	}
	err = sql.Limit(pageSize).
		Offset((page - 1) * pageSize).
		Select("c.*, u.username AS user_name, u.avatar AS user_avatar").
		Find(&comments).Error
	return
}

func (m *ArticleCommentModel) ListReplyComments(ctx context.Context, articleId, commentId int64, page, pageSize int) (comments []ArticleCommentReply, err error) {
	err = m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS c", ArticleComment{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON c.user_id = u.id", User{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u2 ON c.reply_to = u2.id", User{}.TableName())).
		Where("c.parent_id = ?", commentId).
		Where("c.article_id = ?", articleId).
		Where("c.deleted_at IS NULL").
		Order("c.id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Select("c.*, u.username AS user_name, u.avatar AS user_avatar, u2.username AS reply_to_name, u2.avatar AS reply_to_avatar").
		Find(&comments).Error
	return
}

func (m *ArticleCommentModel) GetTotalByArticleId(ctx context.Context, articleId int64) (total int64, err error) {
	return gorm.G[ArticleComment](m.db).Where(generated.ArticleComment.ArticleId.Eq(articleId)).
		Where(generated.ArticleComment.DeletedAt.IsNull()).
		Count(ctx, "*")
}

// 获取指定一级评论id是第几条评论
func (m *ArticleCommentModel) GetCommentIndex(ctx context.Context, commentId uint) (idx int64, err error) {
	err = gorm.G[ArticleComment](m.db).
		Where(generated.ArticleComment.ParentId.Eq(0)).
		Where(generated.ArticleComment.ID.Gt(commentId)).
		Where(generated.ArticleComment.DeletedAt.IsNull()).
		Select("COUNT(*) + 1 AS idx").Scan(ctx, &idx)
	if err != nil {
		return 0, err
	}
	return idx, nil
}

// 获取指定二级评论id是第几条评论
func (m *ArticleCommentModel) GetReplyCommentIndex(ctx context.Context, parentId, commentId uint) (idx int64, err error) {
	err = gorm.G[ArticleComment](m.db).
		Where(generated.ArticleComment.ParentId.Eq(int64(parentId))).
		Where(generated.ArticleComment.ID.Gt(commentId)).
		Where(generated.ArticleComment.DeletedAt.IsNull()).
		Select("COUNT(*) + 1 AS idx").Scan(ctx, &idx)
	if err != nil {
		return 0, err
	}
	return idx, nil
}

// 计算评论回复数量
func (m *ArticleCommentModel) GetReplyCount(ctx context.Context, commentId uint) (count int64, err error) {
	return gorm.G[ArticleComment](m.db).Where(generated.ArticleComment.ParentId.Eq(int64(commentId))).
		Where(generated.ArticleComment.DeletedAt.IsNull()).
		Count(ctx, "*")
}

func (m *ArticleCommentModel) CountCommentByUserId(ctx context.Context, userId uint) (total int64, err error) {
	subSql := gorm.G[Article](m.db).Select("id").Where(generated.Article.AuthorId.Eq(int64(userId)))
	return gorm.G[ArticleComment](m.db).Where("article_id IN (?)", subSql).
		Where(generated.ArticleComment.DeletedAt.IsNull()).
		Count(ctx, "*")
}

func (m *ArticleCommentModel) IncCommentReplyCount(ctx context.Context, commentId uint, count int64) error {
	_, err := gorm.G[ArticleComment](m.db).Where(generated.ArticleComment.ID.Eq(commentId)).
		Set(generated.ArticleComment.ReplyCount.Incr(count)).Update(ctx)
	return err
}

func (m *ArticleCommentModel) PaginationUnreadComments(ctx context.Context, uid int64, page, pageSize int) (comments []ArticleCommentReply, err error) {
	err = m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS c", ArticleComment{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON c.user_id = u.id", User{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u2 ON c.reply_to = u2.id", User{}.TableName())).
		Where("c.reply_to = ?", uid).
		Where("c.deleted_at IS NULL").
		Order("c.id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Select("c.*, u.username AS user_name, u.avatar AS user_avatar, u2.username AS reply_to_name, u2.avatar AS reply_to_avatar").
		Find(&comments).Error
	return
}

// 未读评论总数
func (m *ArticleCommentModel) CountUnreadComments(ctx context.Context, uid int64) (count int64, err error) {
	return gorm.G[ArticleComment](m.db).Where(generated.ArticleComment.ReplyTo.Eq(uid), generated.ArticleComment.DeletedAt.IsNull()).Count(ctx, "id")
}
