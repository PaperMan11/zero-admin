package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
)

// 菜单
func (m *MysqlDB) GetMenus(ctx context.Context, status int32, page, pageSize int) ([]*model.SysMenu, error) {
	if status == 2 {
		return m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.DelFlag.Eq(0)).Order(m.q.SysMenu.Sort.Desc()).Limit(pageSize).Offset((page - 1) * pageSize).Find()
	}
	return m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.Status.Eq(status), m.q.SysMenu.DelFlag.Eq(0)).Order(m.q.SysMenu.Sort.Desc()).Limit(pageSize).Offset((page - 1) * pageSize).Find()
}

func (m *MysqlDB) GetAllMenus(ctx context.Context) ([]*model.SysMenu, error) {
	return m.q.SysMenu.WithContext(ctx).Find()
}

// 根据id获取菜单
func (m *MysqlDB) GetMenuByID(ctx context.Context, menuID int64) (*model.SysMenu, error) {
	return m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.DelFlag.Eq(0), m.q.SysMenu.ID.Eq(menuID)).First()
}

// 根据角色获取有权限的菜单
func (m *MysqlDB) GetMenusByRoles(ctx context.Context, roleCodes []string) ([]*model.SysMenu, error) {
	var scopeIDs []int64
	menu := m.q.SysMenu
	rc := m.q.SysRoleScope
	s := m.q.SysScope
	subQuery := rc.WithContext(ctx).Where(rc.RoleCode.In(roleCodes...))
	if err := s.WithContext(ctx).Where(s.DelFlag.Eq(0), m.q.SysScope.Status.Eq(1), s.Columns(s.ScopeCode).In(subQuery.Select(rc.ScopeCode))).Pluck(s.ID, &scopeIDs); err != nil {
		return nil, err
	}
	if len(scopeIDs) == 0 {
		return []*model.SysMenu{}, nil
	}
	return menu.WithContext(ctx).Where(menu.ScopeID.In(scopeIDs...)).Find()
}

func (m *MysqlDB) GetMenusByScopeIDs(ctx context.Context, scopeIDs []int64) ([]*model.SysMenu, error) {
	return m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.DelFlag.Eq(0), m.q.SysMenu.ScopeID.In(scopeIDs...)).Find()
}

// 创建菜单
func (m *MysqlDB) CreateMenu(ctx context.Context, menu model.SysMenu) (int64, error) {
	err := m.q.SysMenu.WithContext(ctx).Create(&menu)
	return menu.ID, err
}

func (m *MysqlDB) CreateMenus(ctx context.Context, menu []*model.SysMenu) error {
	return m.q.SysMenu.WithContext(ctx).Create(menu...)
}

// 删除菜单
func (m *MysqlDB) DeleteMenu(ctx context.Context, menuID int64) error {
	_, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.ID.Eq(menuID)).Update(m.q.SysMenu.DelFlag, 1)
	return err
}

// 修改菜单
func (m *MysqlDB) UpdateMenu(ctx context.Context, menuID int64, updates interface{}) error {
	_, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.ID.Eq(menuID), m.q.SysMenu.DelFlag.Eq(0)).Updates(updates)
	return err
}

func (m *MysqlDB) SaveMenu(ctx context.Context, menu model.SysMenu) error {
	return m.q.SysMenu.WithContext(ctx).Save(&menu)
}

func (m *MysqlDB) GetMenusByScopeID(ctx context.Context, scopeID int64) ([]*model.SysMenu, error) {
	return m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.DelFlag.Eq(0), m.q.SysMenu.ScopeID.Eq(scopeID)).Find()
}

func (m *MysqlDB) ExistsMenuByName(ctx context.Context, menuName string) (bool, error) {
	count, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.MenuName.Eq(menuName)).Count()
	return count > 0, err
}

func (m *MysqlDB) ExistsMenuByPath(ctx context.Context, menuPath string) (bool, error) {
	count, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.Path.Eq(menuPath)).Count()
	return count > 0, err
}

func (m *MysqlDB) ExistsMenu(ctx context.Context, menuID int64) (bool, error) {
	count, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.ID.Eq(menuID)).Count()
	return count > 0, err
}
