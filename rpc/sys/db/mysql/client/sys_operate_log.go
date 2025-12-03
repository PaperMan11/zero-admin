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

// 获取操作日志详情
func (m *MysqlDB) GetOperateLog(ctx context.Context, logID int64) (*model.SysOperateLog, error) {
	return m.q.SysOperateLog.WithContext(ctx).Where(m.q.SysOperateLog.ID.Eq(logID)).First()
}

// 获取操作日志列表
func (m *MysqlDB) GetOperateLogs(ctx context.Context, filter model.OperateLogFilter, page int, pageSize int) ([]*model.SysOperateLog, int64, error) {
	sql := m.q.SysOperateLog.WithContext(ctx)
	if len(filter.OperationName) > 0 {
		sql.Where(m.q.SysOperateLog.OperationName.Like(filter.OperationName))
	}
	if len(filter.OperationType) > 0 {
		sql.Where(m.q.SysOperateLog.OperationType.Eq(filter.OperationType))
	}
	if len(filter.Title) > 0 {
		sql.Where(m.q.SysOperateLog.Title.Eq(filter.Title))
	}
	if len(filter.OperationIP) > 0 {
		sql.Where(m.q.SysOperateLog.OperationIP.Eq(filter.OperationIP))
	}
	if filter.OperationStatus != 0 {
		sql.Where(m.q.SysOperateLog.OperationStatus.Eq(filter.OperationStatus))
	}
	if len(filter.Browser) > 0 {
		sql.Where(m.q.SysOperateLog.Browser.Eq(filter.Browser))
	}
	if len(filter.Os) > 0 {
		sql.Where(m.q.SysOperateLog.Os.Eq(filter.Os))
	}
	if len(filter.OperationURL) > 0 {
		sql.Where(m.q.SysOperateLog.OperationURL.Eq(filter.OperationURL))
	}
	if len(filter.RequestMethod) > 0 {
		sql.Where(m.q.SysOperateLog.RequestMethod.Eq(filter.RequestMethod))
	}
	count, _ := sql.Count()
	if count == 0 {
		return []*model.SysOperateLog{}, 0, nil
	}
	logs, err := sql.Offset((page - 1) * pageSize).Limit(pageSize).Find()
	return logs, count, err
}

// 删除操作日志
func (m *MysqlDB) DeleteOperateLogs(ctx context.Context, logIDs []int64) error {
	_, err := m.q.SysOperateLog.WithContext(ctx).Where(m.q.SysOperateLog.ID.In(logIDs...)).Delete()
	return err
}

func (m *MysqlDB) DeleteOperateLog(ctx context.Context, logID int64) error {
	_, err := m.q.SysOperateLog.WithContext(ctx).Where(m.q.SysOperateLog.ID.Eq(logID)).Delete()
	return err
}
