package models

import (
	"context"
	"errors"
	"zero-admin/api/blog/internal/models/generated"
	"zero-admin/pkg/errorx"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleCategory struct {
	gorm.Model
	Uid          uint   `json:"uid" gorm:"column:uid"` // 用户ID
	Name         string `json:"name" gorm:"column:name"`
	Description  string `json:"description" gorm:"column:description"`
	ArticleCount int64  `json:"article_count" gorm:"column:article_count"`
}

func (ArticleCategory) TableName() string {
	return "blog_article_category"
}

type ArticleCategoryModel struct {
	db *gorm.DB
}

func NewArticleCategoryModel(db *gorm.DB) *ArticleCategoryModel {
	return &ArticleCategoryModel{db}
}

func (m *ArticleCategoryModel) GetByID(ctx context.Context, id uint) (category *ArticleCategory, err error) {
	var ac ArticleCategory
	ac, err = generated.Query[ArticleCategory](m.db).GetByID(ctx, int64(id))
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errorx.ErrorCategoryNotFound
	case err != nil:
		return nil, errorx.ErrorInternal
	default:
	}
	if ac.DeletedAt.Valid {
		return nil, errorx.ErrorCategoryNotFound
	}
	return &ac, nil
}

func (m *ArticleCategoryModel) PaginationByUid(ctx context.Context, uid uint, page, pageSize int, all bool) (categories []ArticleCategory, err error) {
	sql := gorm.G[ArticleCategory](m.db).
		Where(generated.ArticleCategory.DeletedAt.IsNull()).
		Where(generated.ArticleCategory.Uid.Eq(uid))
	if !all {
		sql.Limit(pageSize).Offset(pageSize * (page - 1))
	}
	return sql.Find(ctx)
}

func (m *ArticleCategoryModel) ListAllByUid(ctx context.Context, uid uint) (categories []ArticleCategory, err error) {
	return gorm.G[ArticleCategory](m.db).
		Where(generated.ArticleCategory.Uid.Eq(uid)).
		Where(generated.ArticleCategory.DeletedAt.IsNull()).
		Order(generated.ArticleCategory.ID.Desc()).
		Find(ctx)
}

// 分类名是否存在
func (m *ArticleCategoryModel) CheckUserCategoryNameExist(ctx context.Context, uid uint, name string) (bool, error) {
	count, err := gorm.G[ArticleCategory](m.db).
		Where(generated.ArticleCategory.Uid.Eq(uid)).
		Where(generated.ArticleCategory.Name.Eq(name)).
		Count(ctx, "*")
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *ArticleCategoryModel) Create(ctx context.Context, category ArticleCategory) (id uint, err error) {
	err = gorm.G[ArticleCategory](m.db).Create(ctx, &category)
	return category.ID, err
}

func (m *ArticleCategoryModel) Update(ctx context.Context, category ArticleCategory) (err error) {
	set := make([]clause.Assigner, 0)
	if category.Name != "" {
		set = append(set, generated.ArticleCategory.Name.Set(category.Name))
	}
	if category.Description != "" {
		set = append(set, generated.ArticleCategory.Description.Set(category.Description))
	}
	_, err = gorm.G[ArticleCategory](m.db).Where(generated.ArticleCategory.ID.Eq(category.ID)).Set(set...).Update(ctx)
	return err
}

func (m *ArticleCategoryModel) Delete(ctx context.Context, id, uid uint) (err error) {
	_, err = gorm.G[ArticleCategory](m.db).Where(generated.ArticleCategory.ID.Eq(id)).Where(generated.ArticleCategory.Uid.Eq(uid)).Delete(ctx)
	return err
}

func (m *ArticleCategoryModel) CountArticleCategoryByUid(ctx context.Context, uid uint) (total int64, err error) {
	return gorm.G[ArticleCategory](m.db).Where(generated.ArticleCategory.Uid.Eq(uid)).Count(ctx, "id")
}

func (m *ArticleCategoryModel) IncArticleCount(ctx context.Context, categoryId uint) (err error) {
	_, err = gorm.G[ArticleCategory](m.db).Where(generated.ArticleCategory.ID.Eq(categoryId)).Set(generated.ArticleCategory.ArticleCount.Incr(1)).Update(ctx)
	return
}

func (m *ArticleCategoryModel) DecArticleCount(ctx context.Context, categoryId uint) (err error) {
	_, err = gorm.G[ArticleCategory](m.db).
		Where(generated.ArticleCategory.ID.Eq(categoryId)).
		Where(generated.ArticleCategory.ArticleCount.Gte(1)).
		Set(generated.ArticleCategory.ArticleCount.Decr(1)).Update(ctx)
	return
}
