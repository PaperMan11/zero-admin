package logic

import (
	"time"
	"zero-admin/pkg/convert"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/sysclient"
)

// ----------------------------------------------user----------------------------------------------

func ConvertToRpcUser(user *model.SysUser) *sysclient.User {
	return &sysclient.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Mobile:   user.Mobile,
		RealName: user.RealName,
		Gender:   user.Gender,
		Status:   user.Status,
		Avatar:   user.Avatar,
	}
}

func ConvertToRpcUsers(users []*model.SysUser) []*sysclient.User {
	res := make([]*sysclient.User, 0, len(users))
	for _, u := range users {
		res = append(res, ConvertToRpcUser(u))
	}
	return res
}

// ----------------------------------------------menu----------------------------------------------

func BuildMenuTree(menus []*model.SysMenu, parentID int64) (menuTree []*sysclient.Menu) {
	menuTree = make([]*sysclient.Menu, 0)
	for _, menu := range menus {
		if menu.ParentID == parentID && menu.Status == 1 {
			m := ConvertToRpcMenu(menu)
			menuTree = append(menuTree, m)
			m.Children = BuildMenuTree(menus, menu.ID)
		}
	}
	return
}

func ConvertToRpcMenu(menu *model.SysMenu) *sysclient.Menu {
	m := &sysclient.Menu{
		Id:        menu.ID,
		ParentId:  menu.ParentID,
		MenuName:  menu.MenuName,
		MenuType:  menu.MenuType,
		Status:    menu.Status,
		ScopeId:   menu.ScopeID,
		Path:      menu.Path,
		Component: menu.Component,
		Redirect:  menu.Redirect,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
	}
	if menu.NoCache == 1 {
		m.NoCache = true
	}
	if menu.Affix == 1 {
		m.Affix = true
	}
	if menu.External == 1 {
		m.External = true
	}
	if menu.Hidden == 1 {
		m.Hidden = true
	}
	return m
}

func ConvertToRpcMenus(menus []*model.SysMenu) []*sysclient.Menu {
	res := make([]*sysclient.Menu, 0, len(menus))
	for _, m := range menus {
		res = append(res, ConvertToRpcMenu(m))
	}
	return res
}

func ConvertToModelMenu(operatorID int64, menu *sysclient.Menu) *model.SysMenu {
	now := time.Now()
	operator := convert.ToString(operatorID)
	m := &model.SysMenu{
		ScopeID:    menu.ScopeId,
		ParentID:   menu.ParentId,
		MenuName:   menu.MenuName,
		MenuType:   menu.MenuType,
		Path:       menu.Path,
		Component:  menu.Component,
		Redirect:   menu.Redirect,
		Icon:       menu.Icon,
		Sort:       menu.Sort,
		Creator:    operator,
		CreateTime: now,
		Updater:    operator,
		UpdateTime: now,
		DelFlag:    0,
		Remark:     menu.Remark,
	}
	if menu.NoCache {
		m.NoCache = 1
	}
	if menu.Affix {
		m.Affix = 1
	}
	if menu.External {
		m.External = 1
	}
	if menu.Hidden {
		m.Hidden = 1
	}
	if menu.Status == 1 {
		m.Status = 1
	}
	return m
}

// ----------------------------------------------role----------------------------------------------

func ConvertToRpcRole(role *model.SysRole) *sysclient.Role {
	return &sysclient.Role{
		RoleId:   role.ID,
		RoleName: role.RoleName,
		RoleCode: role.RoleCode,
	}
}

func ConvertToRpcRoles(roles []*model.SysRole) []*sysclient.Role {
	res := make([]*sysclient.Role, 0, len(roles))
	for _, role := range roles {
		res = append(res, ConvertToRpcRole(role))
	}
	return res
}

// ----------------------------------------------scope----------------------------------------------

func ConvertToRpcScope(scope *model.SysScope) *sysclient.Scope {
	return &sysclient.Scope{
		Id:          scope.ID,
		ScopeName:   scope.ScopeName,
		ScopeCode:   scope.ScopeCode,
		Description: scope.Description,
		Sort:        scope.Sort,
	}
}

func ConvertToRpcScopes(scopes []*model.SysScope) []*sysclient.Scope {
	res := make([]*sysclient.Scope, 0, len(scopes))
	for _, s := range scopes {
		res = append(res, ConvertToRpcScope(s))
	}
	return res
}

// ----------------------------------------------operate log----------------------------------------------

func ConvertToRpcOperateLog(log *model.SysOperateLog) *sysclient.OperateLog {
	return &sysclient.OperateLog{
		Id:                log.ID,
		Title:             log.Title,
		OperationType:     log.OperationType,
		OperationName:     log.OperationName,
		RequestMethod:     log.RequestMethod,
		OperationUrl:      log.OperationURL,
		OperationParams:   log.OperationParams,
		OperationResponse: log.OperationResponse,
		OperationStatus:   log.OperationStatus,
		UseTime:           log.UseTime,
		Browser:           log.Browser,
		Os:                log.Os,
		OperationIp:       log.OperationIP,
		OperationTime:     log.OperationTime.Format(time.DateTime),
	}
}

func ConvertToRpcOperateLogs(logs []*model.SysOperateLog) []*sysclient.OperateLog {
	res := make([]*sysclient.OperateLog, 0, len(logs))
	for _, log := range logs {
		res = append(res, ConvertToRpcOperateLog(log))
	}
	return res
}
