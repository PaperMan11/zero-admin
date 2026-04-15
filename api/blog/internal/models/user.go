package models

import (
	"context"
	"errors"
	"time"
	"zero-admin/api/blog/internal/models/generated"
	"zero-admin/pkg/errorx"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Username     string    `json:"username" gorm:"column:username"`
	Nickname     string    `json:"nickname" gorm:"column:nickname"`
	Password     string    `json:"password" gorm:"column:password"`
	Salt         string    `json:"salt" gorm:"column:salt"`                 // 盐值
	Avatar       string    `json:"avatar" gorm:"column:avatar"`             // 头像
	Gender       int8      `json:"gender" gorm:"column:gender"`             // 性别 0-未知 1-男 2-女
	Introduction string    `json:"introduction" gorm:"column:introduction"` // 个人介绍
	Birthday     string    `json:"birthday" gorm:"column:birthday"`         // 出生日期
	Mobile       string    `json:"mobile" gorm:"column:mobile"`             // 手机号
	Email        string    `json:"email" gorm:"column:email"`
	Status       int8      `json:"status" gorm:"column:status"`               // 状态 0-正常 1-禁用
	LoginTime    time.Time `json:"login_time" gorm:"column:login_time"`       // 最近登录时间
	LastLoginIP  string    `json:"last_login_ip" gorm:"column:last_login_ip"` // 最近登录IP
	LoginCount   int       `json:"login_count" gorm:"column:login_count"`     // 登录次数
}

func (User) TableName() string {
	return "blog_user"
}

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db}
}

// 判断用户是否注销、是否禁用
func CheckUserStatus(user User) error {
	if user.Status == 1 {
		return errorx.ErrorUserDisabled
	}
	if user.DeletedAt.Valid {
		return errorx.ErrorUserNotFound
	}
	return nil
}

func (m *UserModel) GetByID(ctx context.Context, id uint) (user *User, err error) {
	var u User
	u, err = generated.Query[User](m.db).GetByID(ctx, int64(id))
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errorx.ErrorInvalidUsernameOrPassword
	case err != nil:
		return nil, errorx.ErrorInternal
	default:
	}
	if u.Status == 1 {
		return nil, errorx.ErrorUserDisabled
	}
	if u.DeletedAt.Valid {
		return nil, errorx.ErrorUserNotFound
	}
	return &u, nil
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (user User, err error) {
	return generated.Query[User](m.db).FilterWithColumn(ctx, "email", email)
}

func (m *UserModel) GetByMobile(ctx context.Context, mobile string) (user User, err error) {
	return generated.Query[User](m.db).FilterWithColumn(ctx, "mobile", mobile)
}

func (m *UserModel) GetByUsername(ctx context.Context, username string) (user User, err error) {
	return generated.Query[User](m.db).FilterWithColumn(ctx, "username", username)
}

func (m *UserModel) Create(ctx context.Context, user User) (id uint, err error) {
	err = gorm.G[User](m.db).Create(ctx, &user)
	return user.ID, err
}

func (m *UserModel) Save(ctx context.Context, user User) error {
	_, err := gorm.G[User](m.db).Where(generated.User.ID.Eq(user.ID)).Updates(ctx, user)
	return err
}

func (m *UserModel) Updates(ctx context.Context, userId uint, updates map[string]interface{}) (err error) {
	set := clause.Set{}
	for k, v := range updates {
		set = append(set, clause.Assignment{
			Column: clause.Column{Name: k},
			Value:  v,
		})
	}
	_, err = gorm.G[User](m.db).Where(generated.User.ID.Eq(userId)).Set(set).Update(ctx)
	return err
}

func (m *UserModel) Delete(ctx context.Context, userId uint) (err error) {
	_, err = gorm.G[User](m.db).Where(generated.User.ID.Eq(userId)).Update(ctx, "deleted_at", time.Now())
	return err
}
