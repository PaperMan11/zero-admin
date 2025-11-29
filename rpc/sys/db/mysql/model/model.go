package model

type RoleScopeInfo struct {
	RoleID   int64  `json:"role_id" gorm:"column:role_id"`
	RoleName string `json:"role_name" gorm:"column:role_name"`
	RoleCode string `json:"role_code" gorm:"column:role_code"`
	Perm     int32  `json:"perm" gorm:"column:perm"`
	SysScope
}
