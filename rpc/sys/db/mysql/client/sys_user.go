package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/db/mysql/query"
)

// 添加用户
func (u *MysqlDB) CreateUser(ctx context.Context, user model.SysUser) (int64, error) {
	err := u.q.SysUser.WithContext(ctx).Create(&user)
	return user.ID, err
}

// 根据用户名查询用户
func (m *MysqlDB) GetUserByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	return m.q.SysUser.WithContext(ctx).Where(m.q.SysUser.DelFlag.Eq(0), m.q.SysUser.Username.Eq(username)).First()
}

// 根据用户ID查询用户
func (m *MysqlDB) GetUserByID(ctx context.Context, userID int64) (*model.SysUser, error) {
	return m.q.SysUser.WithContext(ctx).Where(m.q.SysUser.DelFlag.Eq(0), m.q.SysUser.ID.Eq(userID)).First()
}

// 更新用户
func (m *MysqlDB) UpdateUserByID(ctx context.Context, userID int64, updates interface{}) error {
	_, err := m.q.SysUser.WithContext(ctx).Where(m.q.SysUser.ID.Eq(userID)).Updates(updates)
	return err
}

func (m *MysqlDB) DeleteUserTx(ctx context.Context, userID int64) error {
	return m.q.Transaction(func(tx *query.Query) error {
		userModel := tx.SysUser
		userRoleModel := tx.SysUserRole
		_, err := userModel.WithContext(ctx).Where(userModel.ID.Eq(userID)).Update(userModel.DelFlag, 1)
		if err != nil {
			return err
		}
		_, err = userRoleModel.WithContext(ctx).Where(userRoleModel.UserID.Eq(userID)).Delete()
		return err
	})
}

func (m *MysqlDB) GetUsersPagination(ctx context.Context, status int32, page, pageSize int) ([]*model.SysUser, error) {
	if status == 2 {
		return m.q.SysUser.WithContext(ctx).Where(m.q.SysUser.DelFlag.Eq(0)).Order(m.q.SysUser.ID.Desc()).Offset((page - 1) * pageSize).Limit(pageSize).Find()
	}
	return m.q.SysUser.WithContext(ctx).Where(m.q.SysUser.DelFlag.Eq(0), m.q.SysUser.Status.Eq(status)).Order(m.q.SysUser.ID.Desc()).Offset((page - 1) * pageSize).Limit(pageSize).Find()
}

func (m *MysqlDB) CountUsers(ctx context.Context, status int32) (int64, error) {
	return m.q.SysUser.WithContext(ctx).Where(m.q.SysUser.DelFlag.Eq(0), m.q.SysUser.Status.Eq(status)).Count()
}

func (m *MysqlDB) SaveUser(ctx context.Context, user model.SysUser) error {
	return m.q.SysUser.WithContext(ctx).Save(&user)
}
