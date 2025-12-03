package model

type RoleScopeInfo struct {
	RoleID   int64  `json:"role_id" gorm:"column:role_id"`
	RoleName string `json:"role_name" gorm:"column:role_name"`
	RoleCode string `json:"role_code" gorm:"column:role_code"`
	Perm     int32  `json:"perm" gorm:"column:perm"`
	SysScope
}

type OperateLogFilter struct {
	Title           string `gorm:"column:title;not null;comment:系统模块" json:"title"`                       // 系统模块
	OperationType   string `gorm:"column:operation_type;not null;comment:操作类型" json:"operation_type"`     // 操作类型
	OperationName   string `gorm:"column:operation_name;not null;comment:操作人员" json:"operation_name"`     // 操作人员
	RequestMethod   string `gorm:"column:request_method;not null;comment:请求方式" json:"request_method"`     // 请求方式
	OperationURL    string `gorm:"column:operation_url;not null;comment:操作方法" json:"operation_url"`       // 操作方法
	OperationStatus int32  `gorm:"column:operation_status;not null;comment:操作状态" json:"operation_status"` // 操作状态
	Browser         string `gorm:"column:browser;not null;comment:浏览器" json:"browser"`                    // 浏览器
	Os              string `gorm:"column:os;not null;comment:操作系统" json:"os"`                             // 操作系统
	OperationIP     string `gorm:"column:operation_ip;not null;comment:操作地址" json:"operation_ip"`         // 操作地址
}
