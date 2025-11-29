package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/db/mysql/query"
)

// 删除角色关联权限
func (m *MysqlDB) DeleteRoleScopes(ctx context.Context, roleCode string, scopeCodes []string) error {
	_, err := m.q.SysRoleScope.WithContext(ctx).
		Where(query.SysRoleScope.RoleCode.Eq(roleCode)).
		Where(query.SysRoleScope.ScopeCode.In(scopeCodes...)).Delete()
	return err
}

func (m *MysqlDB) UpdateRoleScopesTx(ctx context.Context, roleCode string, roleScopes []model.SysRoleScope) error {
	return m.q.Transaction(func(tx *query.Query) error {
		rc := tx.SysRoleScope
		_, err := rc.WithContext(ctx).Where(rc.RoleCode.Eq(roleCode)).Delete()
		if err != nil {
			return err
		}
		newRoleScopes := make([]*model.SysRoleScope, 0, len(roleScopes))
		for _, roleScope := range roleScopes {
			newRoleScopes = append(newRoleScopes, &model.SysRoleScope{
				RoleCode:  roleCode,
				ScopeCode: roleScope.ScopeCode,
				Perm:      roleScope.Perm,
			})
		}
		return rc.WithContext(ctx).Create(newRoleScopes...)
	})
}

func (m *MysqlDB) AddRoleScopes(ctx context.Context, roleScopes []*model.SysRoleScope) error {
	return m.q.SysRoleScope.WithContext(ctx).Create(roleScopes...)
}

// 获取用户安全范围权限
func (m *MysqlDB) GetRoleScopesPerm(ctx context.Context, roleCode string) (res []model.RoleScopeInfo, err error) {
	rs := m.q.SysRoleScope.As("rs")
	r := m.q.SysRole.As("r")
	s := m.q.SysScope.As("s")
	err = rs.WithContext(ctx).
		LeftJoin(r, r.RoleCode.EqCol(rs.RoleCode)).
		LeftJoin(s, s.ScopeCode.EqCol(rs.ScopeCode)).
		Where(rs.RoleCode.Eq(roleCode)).
		Select(s.ALL, r.ID.As("role_id"), r.RoleCode.As("role_code"), r.RoleName.As("role_name"), rs.Perm.As("perm")).
		Scan(&res)
	return res, err
}

func (m *MysqlDB) GetRolesScopesPerm(ctx context.Context, roleCodes []string) (res []model.RoleScopeInfo, err error) {
	rs := m.q.SysRoleScope.As("rs")
	r := m.q.SysRole.As("r")
	s := m.q.SysScope.As("s")
	err = rs.WithContext(ctx).
		LeftJoin(r, r.RoleCode.EqCol(rs.RoleCode)).
		LeftJoin(s, s.ScopeCode.EqCol(rs.ScopeCode)).
		Where(rs.RoleCode.In(roleCodes...)).
		Select(s.ALL, r.ID.As("role_id"), r.RoleCode.As("role_code"), r.RoleName.As("role_name"), rs.Perm.As("perm")).
		Scan(&res)
	return res, err
}
