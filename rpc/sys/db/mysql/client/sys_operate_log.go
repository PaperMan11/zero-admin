package client

import (
	"context"
	"zero-admin/rpc/sys/db/mysql/model"
)

// 添加操作日志
func (m *MysqlDB) CreateOperationLog(ctx context.Context, log model.SysOperateLog) (int64, error) {
	err := m.q.SysOperateLog.WithContext(ctx).Create(&log)
	return log.ID, err
}

func (m *MysqlDB) CreateOperationLogs(ctx context.Context, logs []*model.SysOperateLog) error {
	return m.q.SysOperateLog.WithContext(ctx).Create(logs...)
}
