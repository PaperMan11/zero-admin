package models

import (
	"context"
	"gorm.io/gorm"
	"zero-admin/api/blog/internal/models/generated"
)

type Notice struct {
	gorm.Model
	Title     string `json:"title" gorm:"column:title"`
	Content   string `json:"content" gorm:"column:content"`
	CreatedBy int64  `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy int64  `json:"updatedBy" gorm:"column:updated_by"`
}

func (Notice) TableName() string {
	return "blog_notice"
}

type NoticeModel struct {
	db *gorm.DB
}

func NewNoticeModel(db *gorm.DB) *NoticeModel {
	return &NoticeModel{db: db}
}

// 获取最新公告
func (m *NoticeModel) GetLatest(ctx context.Context) (notice *Notice, err error) {
	return generated.Query[*Notice](m.db).Last(ctx)
}
