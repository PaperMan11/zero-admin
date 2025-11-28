package db

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-admin/pkg/orm"
	"zero-admin/rpc/sys/db/common"
	"zero-admin/rpc/sys/db/mockdb"
	mysqlCli "zero-admin/rpc/sys/db/mysql/client"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/db/mysql/query"
)

type DB interface {
	// ---------------------用户 & 角色---------------------
	// 添加用户
	CreateUser(ctx context.Context, user model.SysUser) (int64, error)
	// 根据用户名查询用户
	GetUserByUsername(ctx context.Context, username string) (model.SysUser, error)
	// 根据用户ID查询用户
	GetUserByID(ctx context.Context, userID int64) (model.SysUser, error)
	// 更新用户
	UpdateUserByID(ctx context.Context, userID int64, updates interface{}) error
	// 先删除再添加
	AddUserRolesTx(ctx context.Context, userID int64, roleIDs []int64) error
	DeleteUserTx(ctx context.Context, userID int64) error
	GetUsersPagination(ctx context.Context, status int32, page, pageSize int) ([]model.SysUser, error)
	CountUsers(ctx context.Context, status int32) (int64, error)

	// 创建角色
	CreateRole(ctx context.Context, role model.SysRole) (int64, error)
	// 删除角色及关联数据
	DeleteRoleTx(ctx context.Context, roleID int64) error
	// 删除角色关联权限
	DeleteRoleScopes(ctx context.Context, roleID int64, scopeCodes []string) error
	// 根据ID获取角色
	GetRoleByID(ctx context.Context, roleID int64) (model.SysRole, error)
	GetRoleByName(ctx context.Context, roleName string) (model.SysRole, error)
	// 判断角色是否存在
	ExistsRoleByName(ctx context.Context, roleName string) (bool, error)
	ExistsRoleByCode(ctx context.Context, roleCode string) (bool, error)
	ExistsRoleByID(ctx context.Context, roleID int64) (bool, error)
	// 获取用户角色
	GetRolesByUserID(ctx context.Context, userID int64) ([]model.SysRole, error)
	// 分页查询角色
	GetRolesPagination(ctx context.Context, status int32, page, pageSize int) ([]model.SysRole, error)
	// 角色总数量
	CountRoles(ctx context.Context, status int32) (int64, error)
	// 查询角色关联数量
	CountRoleAssociated(ctx context.Context, roleID int64) (int64, error)
	// 获取用户角色code
	GetUserRoleCodes(ctx context.Context, userID int64) ([]string, error)
	SaveRole(ctx context.Context, role model.SysRole) error
	UpdateRoleScopesTx(ctx context.Context, roleCode string, roleScopes []model.SysRoleScope) error
	AddRoleScopes(ctx context.Context, roleScopes []model.SysRoleScope) error
	ToggleRoleStatus(ctx context.Context, roleID int64, status int32) error

	// ---------------------菜单 & 权限---------------------
	// 菜单
	GetMenus(ctx context.Context, status int32, page, pageSize int) ([]model.SysMenu, error)
	// 根据id获取菜单
	GetMenuByID(ctx context.Context, menuID int64) (model.SysMenu, error)
	// 根据角色获取有权限的菜单
	GetMenusByRoles(ctx context.Context, roleCodes []string) ([]model.SysMenu, error)
	// 创建菜单
	CreateMenus(ctx context.Context, menu []model.SysMenu) (int64, error)
	// 删除菜单
	DeleteMenu(ctx context.Context, menuID int64) error
	// 修改菜单
	UpdateMenu(ctx context.Context, menuID int64, updates interface{}) error
	SaveMenu(ctx context.Context, menu model.SysMenu) error
	GetMenusByRoleCode(ctx context.Context, roleCode string) ([]model.SysMenu, error)
	GetMenusByScopeID(ctx context.Context, scopeID int64) ([]model.SysMenu, error)
	ExistsMenuByName(ctx context.Context, menuName string) (bool, error)
	ExistsMenuByPath(ctx context.Context, menuPath string) (bool, error)
	ExistsMenu(ctx context.Context, menuID int64) (bool, error)

	// 安全范围
	CreateScope(ctx context.Context, scope model.SysScope) (int64, error)
	CreateScopeTx(ctx context.Context, scope model.SysScope, menuIDs []int64) (int64, error)
	ExistsScope(ctx context.Context, scopeID int64) (bool, error)
	ExistsScopeByCode(ctx context.Context, scopeCode string) (bool, error)
	CountScopes(ctx context.Context) (int64, error)
	SaveScope(ctx context.Context, scope model.SysScope) error
	GetScopeByID(ctx context.Context, scopeID int64) (model.SysScope, error)
	GetScopesPagination(ctx context.Context, page, pageSize int) ([]model.SysScope, error)
	GetScopesByRoleCode(ctx context.Context, roleCode string) ([]model.SysScope, error)
	GetRoleScopesPerm(ctx context.Context, roleCode string) ([]common.RoleScopeInfo, error)
	GetRolesScopesPerm(ctx context.Context, roleCode []string) ([]common.RoleScopeInfo, error)
	AddScopeMenusTx(ctx context.Context, scopeID int64, menus []int64) error // 先删除再添加
	DeleteScope(ctx context.Context, scopeID int64) error
	DeleteScopeTx(ctx context.Context, scopeID int64) error
	DeleteScopeMenus(ctx context.Context, scopeID int64) error

	// 获取用户安全范围权限
	GetUserPerms(ctx context.Context, userID int64) ([]common.RoleScopeInfo, error)

	// ---------------------登录日志 & 操作日志---------------------
	// 添加登录日志
	CreateLoginLog(ctx context.Context, log model.SysLoginLog) error

	// 添加操作日志
	CreateOperationLog(ctx context.Context, log model.SysOperateLog) error
	CreateOperationLogs(ctx context.Context, logs []model.SysOperateLog) error
}

const (
	DB_MOCK  = "mockdb"
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
	case DB_MOCK:
		return mockdb.NewMockDB()
	}
	return nil, ErrDBTypeNotSupport
}
