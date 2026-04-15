package models

import (
	"context"
	"errors"
	"fmt"
	"time"
	"zero-admin/api/blog/internal/models/generated"
	"zero-admin/pkg/errorx"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Article struct {
	gorm.Model
	Title        string `json:"title" gorm:"column:title"`
	Summary      string `json:"summary" gorm:"column:summary"`
	Content      string `json:"content" gorm:"column:content"`
	Tags         string `json:"tags" gorm:"column:tags"`
	CategoryId   int64  `json:"categoryId" gorm:"column:category_id"`
	AuthorId     int64  `json:"authorId" gorm:"column:author_id"`
	ViewCount    int64  `json:"viewCount" gorm:"column:view_count"`
	CommentCount int64  `json:"commentCount" gorm:"column:comment_count"`
	LikeCount    int64  `json:"likeCount" gorm:"column:like_count"`
	Cover        string `json:"cover" gorm:"column:cover"`
}

func (Article) TableName() string {
	return "blog_article"
}

type ArticleDetail struct {
	Article
	AuthorName string `json:"authorName" gorm:"column:author_name"`
}

func (ArticleDetail) TableName() string {
	return "blog_article"
}

type ArticleModel struct {
	db *gorm.DB
}

func NewArticleModel(db *gorm.DB) *ArticleModel {
	return &ArticleModel{db}
}

// GetByID 获取文章详情
func (m *ArticleModel) GetByID(ctx context.Context, id uint) (article *Article, err error) {
	var a Article
	a, err = generated.Query[Article](m.db).GetByID(ctx, int64(id))
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errorx.ErrorArticleNotFound
	case err != nil:
		return nil, errorx.ErrorInternal
	default:
	}
	if a.DeletedAt.Valid {
		return nil, errorx.ErrorArticleNotFound
	}
	return &a, nil
}

func (m *ArticleModel) GetDetailByID(ctx context.Context, id uint) (article *ArticleDetail, err error) {
	err = m.db.Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Where("a.id=?", id).
		Select("a.*, u.username AS author_name").
		First(&article).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errorx.ErrorArticleNotFound
	case err != nil:
		return nil, errorx.ErrorInternal
	default:
	}
	return
}

func (m *ArticleModel) ListArticleByTags(ctx context.Context, tags string, page, pageSize int) (articles []ArticleDetail, err error) {
	err = m.db.Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Where("a.tags LIKE ?", "%"+tags+"%").
		Select("a.*, u.username AS author_name").
		Limit(pageSize).
		Offset(pageSize * (page - 1)).
		Find(&articles).Error
	return
}

// 根据用户ID和分类查询文章列表
func (m *ArticleModel) GetByUidAndCategory(ctx context.Context, uid uint, categoryId int64, page, pageSize int) (articles []ArticleDetail, err error) {
	err = m.db.Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Where("a.category_id = ?", categoryId).
		Where("a.author_id = ?", uid).
		Omit("a.content").
		Select("a.*, u.username AS author_name").
		Limit(pageSize).
		Offset(pageSize * (page - 1)).
		Find(&articles).Error
	return
}

// 查询分类文章总数量
func (m *ArticleModel) GetTotalByCategoryId(ctx context.Context, categoryId int64) (total int64, err error) {
	// err = m.db.Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
	// 	Where("a.category_id = ?", categoryId).
	// 	Count(&total).Error
	return gorm.G[Article](m.db).Where(generated.Article.CategoryId.Eq(categoryId)).Count(ctx, "*")
}

// 获取文章总数
func (m *ArticleModel) GetTotal(ctx context.Context) (total int64, err error) {
	return gorm.G[Article](m.db).Count(ctx, "id")
}

// Create 创建文章
func (m *ArticleModel) Create(ctx context.Context, article Article) (id uint, err error) {
	err = gorm.G[Article](m.db).Create(ctx, &article)
	return article.ID, err
}

// Delete 删除文章
func (m *ArticleModel) Delete(ctx context.Context, id uint) (err error) {
	_, err = gorm.G[Article](m.db).Where(generated.Article.ID.Eq(id)).Update(ctx, "deleted_at", time.Now())
	return err
}

func (m *ArticleModel) Updates(ctx context.Context, id uint, updates map[string]interface{}) (err error) {
	set := clause.Set{}
	for k, v := range updates {
		set = append(set, clause.Assignment{
			Column: clause.Column{Name: k},
			Value:  v,
		})
	}
	_, err = gorm.G[Article](m.db).Where(generated.Article.ID.Eq(id)).Set(set).Update(ctx)
	return err
}

func (m *ArticleModel) CountArticleByUid(ctx context.Context, uid uint) (total int64, err error) {
	return gorm.G[Article](m.db).Where(generated.Article.AuthorId.Eq(int64(uid))).Count(ctx, "id")
}

func (m *ArticleModel) CountArticleByCategoryId(ctx context.Context, uid uint, categoryId int64) (total int64, err error) {
	return gorm.G[Article](m.db).Where(generated.Article.AuthorId.Eq(int64(uid)), generated.Article.CategoryId.Eq(categoryId)).Count(ctx, "id")
}

func (m *ArticleModel) ListRecentArticle(ctx context.Context, page, pageSize int) (articles []ArticleDetail, err error) {
	err = m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Order("a.id DESC").Limit(pageSize).Offset(pageSize * (page - 1)).
		Omit("a.content").
		Select("a.*, u.username AS author_name").Find(&articles).Error
	return
}

func (m *ArticleModel) ListHotArticle(ctx context.Context, page, pageSize int) (articles []ArticleDetail, err error) {
	err = m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Order("a.view_count DESC, a.like_count DESC, a.comment_count DESC").
		Limit(pageSize).Offset(pageSize * (page - 1)).
		Omit("a.content").
		Select("a.*, u.username AS author_name").Find(&articles).Error
	return
}

func (m *ArticleModel) ListVoteArticle(ctx context.Context, page, pageSize int) (articles []ArticleDetail, err error) {
	err = m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Order("a.like_count DESC").Limit(pageSize).Offset(pageSize * (page - 1)).
		Select("a.*, u.username AS author_name").Find(&articles).Error
	return
}

// 随机获取5篇文章
func (m *ArticleModel) GetRandomArticle(ctx context.Context) (articles []ArticleDetail, err error) {
	err = m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Order("RAND()").Limit(5).
		Omit("a.content").
		Select("a.*, u.username AS author_name").Find(&articles).Error
	return
}

func (m *ArticleModel) GetTimeline(ctx context.Context, uid int64) (articles []Article, err error) {
	return gorm.G[Article](m.db).Where(generated.Article.AuthorId.Eq(uid)).Order(generated.Article.ID.Desc()).Omit("content").Find(ctx)
}

// 模糊搜索
func (m *ArticleModel) Search(ctx context.Context, keyword string, page, pageSize int) (articles []ArticleDetail, err error) {
	err = m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Where("a.title LIKE ? OR a.tags LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("a.view_count DESC, a.like_count DESC, a.id DESC").
		Limit(pageSize).Offset(pageSize * (page - 1)).
		Omit("a.content").
		Select("a.*, u.username AS author_name").
		Find(&articles).Error
	return
}

func (m *ArticleModel) CountArticleByKeyword(ctx context.Context, keyword string) (total int64, err error) {
	searchPattern := "%" + keyword + "%"
	return gorm.G[Article](m.db).Where(generated.Article.Title.Like(searchPattern)).Or(generated.Article.Tags.Like(searchPattern)).Count(ctx, "id")
}

func (m *ArticleModel) SearchArticlesByKeyword(ctx context.Context, keyword string, limit int) (articles []ArticleDetail, err error) {
	searchPattern := "%" + keyword + "%"
	err = m.db.WithContext(ctx).Table(fmt.Sprintf("%s AS a", Article{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON a.author_id = u.id", User{}.TableName())).
		Where("a.title LIKE ? OR a.tags LIKE ?", searchPattern, searchPattern).
		Where("a.deleted_at IS NULL").
		Order("a.view_count DESC, a.like_count DESC, a.id DESC").
		Omit("a.content").
		Select("a.*, u.username AS author_name").
		Limit(limit).
		Find(&articles).Error
	return
}
