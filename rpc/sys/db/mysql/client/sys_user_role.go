package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/db/mysql/query"
)

// 先删除再添加
func (m *MysqlDB) AddUserRolesTx(ctx context.Context, userID int64, roleCodes []string) error {
	return m.q.Transaction(func(tx *query.Query) error {
		userRoleModel := tx.SysUserRole.WithContext(ctx)
		_, err := userRoleModel.Where(tx.SysUserRole.UserID.Eq(userID)).Delete()
		if err != nil {
			return err
		}
		newUserRoles := make([]*model.SysUserRole, 0, len(roleCodes))
		for _, roleCode := range roleCodes {
			newUserRoles = append(newUserRoles, &model.SysUserRole{
				UserID:   userID,
				RoleCode: roleCode,
			})
		}
		return userRoleModel.Create(newUserRoles...)
	})
}

// 获取用户角色
func (m *MysqlDB) GetRolesByUserID(ctx context.Context, userID int64) ([]*model.SysRole, error) {
	ur := m.q.SysUserRole
	r := m.q.SysRole
	subQuery := ur.WithContext(ctx).Where(ur.UserID.Eq(userID))
	return r.WithContext(ctx).Where(r.Columns(r.RoleCode).In(subQuery.Select(ur.RoleCode))).Find()
}

// 查询角色被用户关联数量
func (m *MysqlDB) CountUserRoles(ctx context.Context, roleCode string) (int64, error) {
	return m.q.SysUserRole.WithContext(ctx).Where(m.q.SysUserRole.RoleCode.Eq(roleCode)).Count()
}

// 获取用户角色code
func (m *MysqlDB) GetUserRoleCodes(ctx context.Context, userID int64) (roleCodes []string, err error) {
	err = m.q.SysUserRole.WithContext(ctx).Where(m.q.SysUserRole.UserID.Eq(userID)).Pluck(m.q.SysUserRole.RoleCode, &roleCodes)
	return
}
