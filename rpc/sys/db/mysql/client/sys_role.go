package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/db/mysql/query"
)

// 创建角色
func (m *MysqlDB) CreateRole(ctx context.Context, role model.SysRole) (int64, error) {
	err := m.q.SysRole.WithContext(ctx).Create(&role)
	return role.ID, err
}

// 删除角色及关联数据
func (m *MysqlDB) DeleteRoleTx(ctx context.Context, roleCode string) error {
	return m.q.Transaction(func(tx *query.Query) error {
		_, err := tx.SysRole.WithContext(ctx).Where(tx.SysRole.RoleCode.Eq(roleCode)).Delete()
		if err != nil {
			return err
		}
		_, err = tx.SysUserRole.WithContext(ctx).Where(tx.SysUserRole.RoleCode.Eq(roleCode)).Delete()
		return err
	})
}

// 根据ID获取角色
func (m *MysqlDB) GetRoleByID(ctx context.Context, roleID int64) (*model.SysRole, error) {
	return m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.ID.Eq(roleID)).First()
}

func (m *MysqlDB) GetRoleByIDs(ctx context.Context, roleIDs []int64) ([]*model.SysRole, error) {
	return m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.ID.In(roleIDs...)).Find()
}

func (m *MysqlDB) GetRoleByCode(ctx context.Context, roleCode string) (*model.SysRole, error) {
	return m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.RoleCode.Eq(roleCode)).First()
}

func (m *MysqlDB) GetRoleByCodes(ctx context.Context, roleCodes []string) ([]*model.SysRole, error) {
	return m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.RoleCode.In(roleCodes...)).Find()
}

func (m *MysqlDB) GetRoleByName(ctx context.Context, roleName string) (*model.SysRole, error) {
	return m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.RoleName.Eq(roleName)).First()
}

// 判断角色是否存在
func (m *MysqlDB) ExistsRoleByName(ctx context.Context, roleName string) (bool, error) {
	count, err := m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.RoleName.Eq(roleName)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (m *MysqlDB) ExistsRoleByCode(ctx context.Context, roleCode string) (bool, error) {
	count, err := m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.RoleCode.Eq(roleCode)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (m *MysqlDB) ExistsRoleByID(ctx context.Context, roleID int64) (bool, error) {
	count, err := m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.ID.Eq(roleID)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, err
}

// 分页查询角色
func (m *MysqlDB) GetRolesPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysRole, error) {
	return m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.Status.Eq(status)).Order(m.q.SysRole.ID.Desc()).Limit(pageSize).Offset((page - 1) * pageSize).Find()
}

// 角色总数量
func (m *MysqlDB) CountRoles(ctx context.Context, status int32) (int64, error) {
	return m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.Status.Eq(status)).Count()
}

func (m *MysqlDB) SaveRole(ctx context.Context, role model.SysRole) error {
	return m.q.SysRole.WithContext(ctx).Save(&role)
}

func (m *MysqlDB) ToggleRoleStatus(ctx context.Context, roleID int64, status int32, operator string) error {
	_, err := m.q.SysRole.WithContext(ctx).Where(m.q.SysRole.ID.Eq(roleID)).Updates(map[string]interface{}{
		"status":  status,
		"updater": operator,
	})
	return err
}
