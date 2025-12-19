package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/db/mysql/query"
)

func (m *MysqlDB) CreateScope(ctx context.Context, scope model.SysScope) (int64, error) {
	err := m.q.SysScope.WithContext(ctx).Create(&scope)
	return scope.ID, err
}

func (m *MysqlDB) CreateScopeTx(ctx context.Context, scope model.SysScope, menuIDs []int64) (int64, error) {
	err := m.q.Transaction(func(tx *query.Query) error {
		err := tx.SysScope.WithContext(ctx).Create(&scope)
		if err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		_, err = tx.SysMenu.WithContext(ctx).Where(tx.SysMenu.ID.In(menuIDs...)).Update(tx.SysMenu.ScopeID, scope.ID)
		return err
	})
	return scope.ID, err
}

func (m *MysqlDB) ExistsScope(ctx context.Context, scopeID int64) (bool, error) {
	count, err := m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.ID.Eq(scopeID)).Count()
	return count > 0, err
}

func (m *MysqlDB) ExistsScopeByCode(ctx context.Context, scopeCode string) (bool, error) {
	count, err := m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.ScopeCode.Eq(scopeCode)).Count()
	return count > 0, err
}

func (m *MysqlDB) CountScopes(ctx context.Context) (int64, error) {
	count, err := m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.DelFlag.Eq(0)).Count()
	return count, err
}

func (m *MysqlDB) SaveScope(ctx context.Context, scope model.SysScope) error {
	return m.q.SysScope.WithContext(ctx).Save(&scope)
}

func (m *MysqlDB) GetScopeByID(ctx context.Context, scopeID int64) (*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.DelFlag.Eq(0), m.q.SysScope.ID.Eq(scopeID)).First()
}

func (m *MysqlDB) GetScopeByCode(ctx context.Context, scopeCode string) (*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.DelFlag.Eq(0), m.q.SysScope.ScopeCode.Eq(scopeCode)).First()
}

func (m *MysqlDB) GetScopes(ctx context.Context, scopeIDs []int64) ([]*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.DelFlag.Eq(0), m.q.SysScope.ID.In(scopeIDs...)).Find()
}

func (m *MysqlDB) GetAllScopes(ctx context.Context) ([]*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.DelFlag.Eq(0)).Find()
}

func (m *MysqlDB) GetScopesByCodes(ctx context.Context, scopeCodes []string) ([]*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.DelFlag.Eq(0), m.q.SysScope.ScopeCode.In(scopeCodes...)).Find()
}

func (m *MysqlDB) GetScopesPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysScope, error) {
	if status == 2 {
		return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.DelFlag.Eq(0)).Order(m.q.SysScope.Sort.Desc()).Offset((page - 1) * pageSize).Limit(pageSize).Find()
	}
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.DelFlag.Eq(0), m.q.SysScope.Status.Eq(status)).Order(m.q.SysScope.Sort.Desc()).Offset((page - 1) * pageSize).Limit(pageSize).Find()
}

func (m *MysqlDB) AddScopeMenus(ctx context.Context, scopeID int64, menus []int64) error {
	_, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.ID.In(menus...)).Update(m.q.SysMenu.ScopeID, scopeID)
	return err
}

func (m *MysqlDB) DeleteScopeTx(ctx context.Context, scopeID int64) error {
	return m.q.Transaction(func(tx *query.Query) error {
		_, err := tx.SysMenu.WithContext(ctx).Where(tx.SysMenu.ScopeID.Eq(scopeID)).Update(tx.SysMenu.ScopeID, 0)
		if err != nil {
			return err
		}
		_, err = tx.SysScope.WithContext(ctx).Where(m.q.SysScope.ID.Eq(scopeID)).Update(tx.SysScope.DelFlag, 1)
		return err
	})
}

func (m *MysqlDB) DeleteScopeMenus(ctx context.Context, scopeID int64) error {
	_, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.ScopeID.Eq(scopeID)).Update(m.q.SysMenu.ScopeID, 0)
	return err
}

func (m *MysqlDB) ToggleScopeStatus(ctx context.Context, scopeID int64, status int32, operator string) error {
	_, err := m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.ID.Eq(scopeID)).Updates(map[string]interface{}{
		"status":  status,
		"updater": operator,
	})
	return err
}

func (m *MysqlDB) GetUnassignedMenus(ctx context.Context) ([]*model.SysMenu, error) {
	return m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.DelFlag.Eq(0), m.q.SysMenu.ScopeID.Eq(0)).Find()
}

// 全量更新安全范围的菜单树
func (m *MysqlDB) UpdateScopeMenusTx(ctx context.Context, scopeID int64, menus []*model.SysMenu) error {
	return m.q.Transaction(func(tx *query.Query) error {
		_, err := tx.SysMenu.WithContext(ctx).Where(tx.SysMenu.ScopeID.Eq(scopeID)).Update(tx.SysMenu.ScopeID, 0)
		if err != nil {
			return err
		}

		if len(menus) > 0 {
			for _, menu := range menus {
				tx.SysMenu.WithContext(ctx).Where(tx.SysMenu.ID.Eq(menu.ID)).Updates(map[string]interface{}{
					"scope_id":  scopeID,
					"parent_id": menu.ParentID,
				})
			}
		}

		return nil
	})
}
