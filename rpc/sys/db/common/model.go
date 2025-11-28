package common

import "zero-admin/rpc/sys/db/mysql/model"

type RoleScopeInfo struct {
	RoleID   int64
	RoleName string
	RoleCode string
	Perm     int32
	model.SysScope
}
