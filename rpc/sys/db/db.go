package db

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-admin/pkg/orm"
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
	GetUserByUsername(ctx context.Context, username string) (*model.SysUser, error)
	// 根据用户ID查询用户
	GetUserByID(ctx context.Context, userID int64) (*model.SysUser, error)
	// 更新用户
	UpdateUserByID(ctx context.Context, userID int64, updates interface{}) error
	// 先删除再添加
	AddUserRolesTx(ctx context.Context, userID int64, roleCodes []string) error
	DeleteUserTx(ctx context.Context, userID int64) error
	GetUsersPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysUser, error)
	CountUsers(ctx context.Context, status int32) (int64, error)
	SaveUser(ctx context.Context, user model.SysUser) error

	// 创建角色
	CreateRole(ctx context.Context, role model.SysRole) (int64, error)
	// 删除角色及关联数据
	DeleteRoleTx(ctx context.Context, roleCode string) error
	// 删除角色关联权限
	DeleteRoleScopes(ctx context.Context, roleCode string, scopeCodes []string) error
	// 根据ID获取角色
	GetRoleByID(ctx context.Context, roleID int64) (*model.SysRole, error)
	GetRoleByIDs(ctx context.Context, roleIDs []int64) ([]*model.SysRole, error)
	GetRoleByName(ctx context.Context, roleName string) (*model.SysRole, error)
	GetRoleByCode(ctx context.Context, roleCode string) (*model.SysRole, error)
	GetRoleByCodes(ctx context.Context, roleCodes []string) ([]*model.SysRole, error)
	// 判断角色是否存在
	ExistsRoleByName(ctx context.Context, roleName string) (bool, error)
	ExistsRoleByCode(ctx context.Context, roleCode string) (bool, error)
	ExistsRoleByID(ctx context.Context, roleID int64) (bool, error)
	// 获取用户角色
	GetRolesByUserID(ctx context.Context, userID int64) ([]*model.SysRole, error)
	// 分页查询角色
	GetRolesPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysRole, error)
	GetAllRoles(ctx context.Context) ([]*model.SysRole, error)
	// 角色总数量
	CountRoles(ctx context.Context, status int32) (int64, error)
	// 查询角色被用户关联数量
	CountUserRoles(ctx context.Context, roleCode string) (int64, error)
	// 获取用户角色code
	GetUserRoleCodes(ctx context.Context, userID int64) ([]string, error)
	SaveRole(ctx context.Context, role model.SysRole) error
	UpdateRoleScopesTx(ctx context.Context, roleCode string, roleScopes []model.SysRoleScope) error
	AddRoleScopes(ctx context.Context, roleScopes []*model.SysRoleScope) error
	UpsertRoleScopes(ctx context.Context, roleScope model.SysRoleScope) error
	ToggleRoleStatus(ctx context.Context, roleID int64, status int32, operator string) error
	GetRolesByScopeCode(ctx context.Context, scopeCode string) ([]*model.SysRole, error)
	GetRolePermsByScopeCode(ctx context.Context, scopeCode string) ([]*model.SysRoleScope, error)

	// ---------------------菜单 & 权限---------------------
	// 菜单
	GetMenus(ctx context.Context, status int32, page, pageSize int) ([]*model.SysMenu, error)
	GetAllMenus(ctx context.Context) ([]*model.SysMenu, error)
	GetUnassignedMenus(ctx context.Context) ([]*model.SysMenu, error) // 未分配的菜单
	// 根据id获取菜单
	GetMenuByID(ctx context.Context, menuID int64) (*model.SysMenu, error)
	// 根据角色获取有权限的菜单
	GetMenusByRoles(ctx context.Context, roleCodes []string) ([]*model.SysMenu, error)
	// 创建菜单
	CreateMenu(ctx context.Context, menu model.SysMenu) (int64, error)
	CreateMenus(ctx context.Context, menu []*model.SysMenu) error
	// 删除菜单
	DeleteMenu(ctx context.Context, menuID int64) error
	// 修改菜单
	UpdateMenu(ctx context.Context, menuID int64, updates interface{}) error
	SaveMenu(ctx context.Context, menu model.SysMenu) error
	GetMenusByScopeID(ctx context.Context, scopeID int64) ([]*model.SysMenu, error)
	GetMenusByScopeIDs(ctx context.Context, scopeIDs []int64) ([]*model.SysMenu, error)
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
	GetScopeByID(ctx context.Context, scopeID int64) (*model.SysScope, error)
	GetScopeByCode(ctx context.Context, scopeCode string) (*model.SysScope, error)
	GetScopes(ctx context.Context, scopeIDs []int64) ([]*model.SysScope, error)
	GetAllScopes(ctx context.Context) ([]*model.SysScope, error)
	GetScopesByCodes(ctx context.Context, scopeCodes []string) ([]*model.SysScope, error)
	GetScopesPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysScope, error)
	AddScopeMenus(ctx context.Context, scopeID int64, menus []int64) error // 先删除再添加
	DeleteScopeTx(ctx context.Context, scopeID int64) error
	DeleteScopeMenus(ctx context.Context, scopeID int64) error
	ToggleScopeStatus(ctx context.Context, scopeID int64, status int32, operator string) error
	UpdateScopeMenusTx(ctx context.Context, scopeID int64, menus []*model.SysMenu) error // 全量更新安全范围的菜单树

	// 获取用户安全范围权限
	GetRoleScopesPerm(ctx context.Context, roleCode string) ([]model.RoleScopeInfo, error)
	GetRolesScopesPerm(ctx context.Context, roleCodes []string) ([]model.RoleScopeInfo, error)

	// ---------------------登录日志 & 操作日志---------------------
	// 添加登录日志
	CreateLoginLog(ctx context.Context, log model.SysLoginLog) (int64, error)

	// 添加操作日志
	CreateOperationLog(ctx context.Context, log model.SysOperateLog) (int64, error)
	CreateOperationLogs(ctx context.Context, logs []*model.SysOperateLog) error
	// 获取操作日志详情
	GetOperateLog(ctx context.Context, logID int64) (*model.SysOperateLog, error)
	// 获取操作日志列表
	GetOperateLogs(ctx context.Context, filter model.OperateLogFilter, page int, pageSize int) ([]*model.SysOperateLog, int64, error)
	DeleteOperateLogs(ctx context.Context, logIDs []int64) error // 删除操作日志
	DeleteOperateLog(ctx context.Context, logID int64) error
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
