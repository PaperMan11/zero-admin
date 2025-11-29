package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
)

// 添加登录日志
func (m *MysqlDB) CreateLoginLog(ctx context.Context, log model.SysLoginLog) (int64, error) {
	err := m.q.SysLoginLog.WithContext(ctx).Create(&log)
	return log.ID, err
}
