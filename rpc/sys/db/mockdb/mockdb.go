package mockdb

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
)

type MockDB struct {
}

func NewMockDB() (*MockDB, error) {
	return &MockDB{}, nil
}

// ---------------------用户 & 角色---------------------
// 添加用户
func (*MockDB) CreateUser(ctx context.Context, user model.SysUser) (int64, error) {
	return 0, nil
}

// 根据用户名查询用户
func (*MockDB) GetUserByUsername(ctx context.Context, username string) (model.SysUser, error) {
	return model.SysUser{}, nil
}

// 根据用户ID查询用户
func (*MockDB) GetUserByID(ctx context.Context, userID int64) (model.SysUser, error) {
	return model.SysUser{}, nil
}

// 更新用户
func (*MockDB) UpdateUserByID(ctx context.Context, userID int64, updates interface{}) error {
	return nil
}

// 创建角色
func (*MockDB) CreateRole(ctx context.Context, role model.SysRole) (int64, error) {
	return 0, nil
}

// 根据ID获取角色
func (*MockDB) GetRoleByID(ctx context.Context, roleID int64) (model.SysRole, error) {
	return model.SysRole{}, nil
}

// 获取用户角色
func (*MockDB) GetRolesByUserID(ctx context.Context, userID int64) ([]model.SysRole, error) {
	return nil, nil
}

// ---------------------菜单 & 权限---------------------
// 获取所有的菜单
func (*MockDB) GetMenus(ctx context.Context, page, pageSize int) ([]model.SysMenu, error) {
	return nil, nil
}

// 根据id获取菜单
func (*MockDB) GetMenuByID(ctx context.Context, menuID int64) (model.SysMenu, error) {
	return model.SysMenu{}, nil
}

// 根据角色获取有权限的菜单
func (*MockDB) GetMenusByRoles(ctx context.Context, roleCodes []string) ([]model.SysMenu, error) {
	return nil, nil
}

// 创建菜单
func (*MockDB) CreateMenus(ctx context.Context, menu []model.SysMenu) (int64, error) {
	return 0, nil
}

// 删除菜单
func (*MockDB) DeleteMenu(ctx context.Context, menuID int64) error {
	return nil
}

// 修改菜单
func (*MockDB) UpdateMenu(ctx context.Context, menuID int64, updates interface{}) error {
	return nil
}
func (*MockDB) GetMenusByRoleID(ctx context.Context, roleID int64) ([]model.SysMenu, error) {
	return nil, nil
}
func (*MockDB) GetMenusByScopeID(ctx context.Context, scopeID int64) ([]model.SysMenu, error) {
	return nil, nil
}

// 创建安全范围
func (*MockDB) CreateScope(ctx context.Context, scope model.SysScope) (int64, error) {
	return 0, nil
}

// 获取安全范围
func (*MockDB) GetScopeByID(ctx context.Context, scopeID int64) (model.SysScope, error) {
	return model.SysScope{}, nil
}
func (*MockDB) GetScopesByRoleID(ctx context.Context, roleID int64) ([]model.SysScope, error) {
	return nil, nil
}

// ---------------------登录日志---------------------
// 添加登录日志
func (*MockDB) CreateLoginLog(ctx context.Context, log model.SysLoginLog) error {
	return nil
}
