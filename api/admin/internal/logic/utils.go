package logic

import (
	"context"
	"zero-admin/api/admin/internal/types"
	"zero-admin/pkg/convert"
	"zero-admin/rpc/sys/client/roleservice"
	"zero-admin/rpc/sys/client/scopeservice"
	"zero-admin/rpc/sys/client/userservice"
)

func GetOperateID(ctx context.Context) int64 {
	return convert.ToInt64(ctx.Value("uid"))
}

func ConvertToTypesMenus(menus []*scopeservice.Menu) []types.Menu {
	menuList := make([]types.Menu, 0, len(menus))
	for _, v := range menus {
		menuList = append(menuList, ConvertToTypesMenu(v))
	}
	return menuList
}

func ConvertToTypesMenu(menu *scopeservice.Menu) types.Menu {
	return types.Menu{
		Id:        menu.Id,
		ParentId:  menu.ParentId,
		MenuType:  menu.MenuType,
		Path:      menu.Path,
		Component: menu.Component,
		Redirect:  menu.Redirect,
		Sort:      menu.Sort,
		External:  menu.External,
		Hidden:    menu.Hidden,
		Status:    menu.Status,
		ScopeId:   menu.ScopeId,
		Remark:    menu.Remark,
		Meta: types.MenuMeta{
			MenuName: menu.MenuName,
			Icon:     menu.Icon,
			NoCache:  menu.NoCache,
			Affix:    menu.Affix,
		},
		Children: ConvertToTypesMenus(menu.Children),
		Perms:    menu.Perms,
	}
}

func ConvertToTypesScope(scope *scopeservice.Scope) types.Scope {
	return types.Scope{
		Id:          scope.Id,
		ScopeName:   scope.ScopeName,
		ScopeCode:   scope.ScopeCode,
		Description: scope.Description,
		Sort:        scope.Sort,
	}
}

func ConvertToTypesRole(role *roleservice.Role) types.Role {
	return types.Role{
		RoleId:      role.RoleId,
		RoleName:    role.RoleName,
		RoleCode:    role.RoleCode,
		Description: role.Description,
		Status:      role.Status,
	}
}

func ConvertToTypesRoles(roles []*roleservice.Role) []types.Role {
	r := make([]types.Role, 0, len(roles))
	for _, v := range roles {
		r = append(r, ConvertToTypesRole(v))
	}
	return r
}

func ConvertToTypesUser(user *userservice.User) types.User {
	return types.User{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Mobile:   user.Mobile,
		RealName: user.RealName,
		Gender:   user.Gender,
		Status:   user.Status,
		Avatar:   user.Avatar,
	}
}

func ConvertToTypesUsers(users []*userservice.User) []types.User {
	u := make([]types.User, 0, len(users))
	for _, v := range users {
		u = append(u, ConvertToTypesUser(v))
	}
	return u
}

func ConvertToTypesUserInfo(userInfo *userservice.UserInfo) types.UserInfo {
	return types.UserInfo{
		User: types.User{
			Id:       userInfo.Id,
			Username: userInfo.Username,
			Email:    userInfo.Email,
			Mobile:   userInfo.Mobile,
			RealName: userInfo.RealName,
			Gender:   userInfo.Gender,
			Status:   userInfo.Status,
			Avatar:   userInfo.Avatar,
		},
		Roles:    ConvertToTypesRoles(userInfo.Roles),
		MenuTree: ConvertToTypesMenus(userInfo.MenuTree),
	}
}
