package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
)

// 添加用户
func (m *MysqlDB) CreateUser(ctx context.Context, user model.SysUser) error {
	return m.q.SysUser.WithContext(ctx).Create(&user)
}

// 根据用户名查询用户
func (m *MysqlDB) GetUserByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	u := m.q.SysUser.WithContext(ctx)
	return u.Where(m.q.SysUser.Username.Eq(username)).First()
}
