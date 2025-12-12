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
	count, err := m.q.SysScope.WithContext(ctx).Count()
	return count, err
}

func (m *MysqlDB) SaveScope(ctx context.Context, scope model.SysScope) error {
	return m.q.SysScope.WithContext(ctx).Save(&scope)
}

func (m *MysqlDB) GetScopeByID(ctx context.Context, scopeID int64) (*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.ID.Eq(scopeID)).First()
}

func (m *MysqlDB) GetScopes(ctx context.Context, scopeIDs []int64) ([]*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.ID.In(scopeIDs...)).Find()
}

func (m *MysqlDB) GetAllScopes(ctx context.Context) ([]*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Find()
}

func (m *MysqlDB) GetScopesByCodes(ctx context.Context, scopeCodes []string) ([]*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.ScopeCode.In(scopeCodes...)).Find()
}

func (m *MysqlDB) GetScopesPagination(ctx context.Context, page, pageSize int) ([]*model.SysScope, error) {
	return m.q.SysScope.WithContext(ctx).Order(m.q.SysScope.Sort.Desc()).Offset((page - 1) * pageSize).Limit(pageSize).Find()
}

func (m *MysqlDB) AddScopeMenus(ctx context.Context, scopeID int64, menus []int64) error {
	_, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.ID.In(menus...)).Update(m.q.SysMenu.ScopeID, scopeID)
	return err
}

func (m *MysqlDB) DeleteScope(ctx context.Context, scopeID int64) error {
	_, err := m.q.SysScope.WithContext(ctx).Where(m.q.SysScope.ID.Eq(scopeID)).Delete()
	return err
}

func (m *MysqlDB) DeleteScopeMenus(ctx context.Context, scopeID int64) error {
	_, err := m.q.SysMenu.WithContext(ctx).Where(m.q.SysMenu.ScopeID.Eq(scopeID)).Update(m.q.SysMenu.ScopeID, 0)
	return err
}
