package db

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-admin/pkg/orm"
	mysqlCli "zero-admin/rpc/sys/db/mysql/client"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/db/mysql/query"
)

type UserWithRoleInfo struct {
	model.SysUser
	Roles []model.SysRole
}

type DB interface {
	// ---------------------用户 & 角色---------------------
	// 添加用户
	CreateUser(ctx context.Context, user model.SysUser) error
	// 根据用户名查询用户
	GetUserByUsername(ctx context.Context, username string) (*model.SysUser, error)
	// 获取用户及角色
	GetUserWithRole(ctx context.Context, userID int64) (UserWithRoleInfo, error)
	GetUserWithRoleByUsername(ctx context.Context, username string) (UserWithRoleInfo, error)
	// 更新用户
	UpdateUserByID(ctx context.Context, userID int64, updates interface{}) error

	// ---------------------菜单 & 权限---------------------
	// 获取所有的菜单
	GetMenus(ctx context.Context) ([]model.SysMenu, error)
	// 根据id获取菜单
	GetMenuByID(ctx context.Context, menuID int64) (*model.SysMenu, error)
	// 根据角色获取有权限的菜单
	GetMenusByRole(ctx context.Context, roleCodes []string) ([]model.SysMenu, error)
	// 创建菜单
	CreateMenus(ctx context.Context, menu []model.SysMenu) (int64, error)
	// 删除菜单
	DeleteMenu(ctx context.Context, menuID int64) error
	// 修改菜单
	UpdateMenu(ctx context.Context, menuID int64, updates interface{}) error

	// ---------------------登录日志---------------------
	// 添加登录日志
	CreateLoginLog(ctx context.Context, log model.SysLoginLog) error
}

const (
	DB_MYSQL = "mysql"
)

var (
	ErrDBTypeNotSupport       = errors.New("db type not support")
	ErrDBConfigTypeNotSupport = errors.New("db config type not support")
)

func MustNewDB(dbType string, dbConfig interface{}) DB {
	db, err := NewDB(dbType, dbConfig)
	logx.Must(err)
	return db
}

func NewDB(dbType string, dbConfig interface{}) (DB, error) {
	switch dbType {
	case DB_MYSQL:
		dbConf, ok := dbConfig.(*orm.Config)
		if !ok {
			return nil, ErrDBConfigTypeNotSupport
		}
		dbCli, err := orm.NewMysql(dbConf)
		if err != nil {
			return nil, err
		}
		q := query.Use(dbCli)
		return mysqlCli.NewMysqlDB(q)
	}
	return nil, ErrDBTypeNotSupport
}
