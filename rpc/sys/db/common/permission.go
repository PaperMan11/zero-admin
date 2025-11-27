package common

import "strings"

type PermType = int32

const (
	// 为了方便管理，用户2 4 8权限的都自动加上1权限
	// 0-无权限，1-读，2-写，4-创建，8-删除
	PERM_NONE PermType = 0

	PERM_READ = 1 << (iota - 1)
	PERM_UPDATE
	PERM_CREATE
	PERM_DELETE
	PERM_READ_UPDATE              PermType = PERM_READ | PERM_UPDATE
	PERM_READ_CREATE              PermType = PERM_READ | PERM_CREATE
	PERM_READ_DELETE              PermType = PERM_READ | PERM_DELETE
	PERM_UPDATE_CREATE            PermType = PERM_UPDATE | PERM_CREATE
	PERM_UPDATE_DELETE            PermType = PERM_UPDATE | PERM_DELETE
	PERM_CREATE_DELETE            PermType = PERM_CREATE | PERM_DELETE
	PERM_READ_UPDATE_CREATE       PermType = PERM_READ | PERM_UPDATE | PERM_CREATE
	PERM_READ_UPDATE_DELETE       PermType = PERM_READ | PERM_UPDATE | PERM_DELETE
	PERM_READ_CREATE_DELETE       PermType = PERM_READ | PERM_CREATE | PERM_DELETE
	PERM_UPDATE_CREATE_DELETE     PermType = PERM_UPDATE | PERM_CREATE | PERM_DELETE
	PERM_READ_WRITE_CREATE_DELETE PermType = PERM_READ | PERM_UPDATE | PERM_CREATE | PERM_DELETE
	PERM_ALL                      PermType = PERM_READ_WRITE_CREATE_DELETE
)

var PermissionMap = map[PermType][]string{
	PERM_NONE:   {"NONE"},
	PERM_READ:   {"READ"},
	PERM_CREATE: {"CREATE"},
	PERM_UPDATE: {"UPDATE"},
	PERM_DELETE: {"DELETE"},

	PERM_READ_CREATE:   {"READ", "CREATE"},
	PERM_READ_UPDATE:   {"READ", "UPDATE"},
	PERM_READ_DELETE:   {"READ", "DELETE"},
	PERM_UPDATE_CREATE: {"CREATE", "UPDATE"},
	PERM_CREATE_DELETE: {"CREATE", "DELETE"},
	PERM_UPDATE_DELETE: {"UPDATE", "DELETE"},

	PERM_READ_UPDATE_CREATE:   {"READ", "CREATE", "UPDATE"},
	PERM_READ_CREATE_DELETE:   {"READ", "CREATE", "DELETE"},
	PERM_READ_UPDATE_DELETE:   {"READ", "UPDATE", "DELETE"},
	PERM_UPDATE_CREATE_DELETE: {"CREATE", "UPDATE", "DELETE"},

	PERM_READ_WRITE_CREATE_DELETE: {"READ", "CREATE", "UPDATE", "DELETE"},
}

func ParsePermission(perms []string) (p PermType) {
	if len(perms) == 0 {
		return PERM_NONE
	}
	p = PERM_READ // 为方便处理默认都有读权限
	for _, perm := range perms {
		switch strings.ToUpper(perm) {
		case "READ":
			p |= PERM_READ
		case "UPDATE":
			p |= PERM_UPDATE
		case "DELETE":
			p |= PERM_DELETE
		case "CREATE":
			p |= PERM_CREATE
		}
	}
	return p
}
