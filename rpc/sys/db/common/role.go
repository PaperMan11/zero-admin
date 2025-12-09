package common

import "errors"

const (
	SUPERUSER = "SUPERUSER" // 超级用户
)

var ErrSuperUserDoNotEdit = errors.New("do not edit")

func IsSuperUser(roleName string) bool {
	return roleName == SUPERUSER
}
