package authservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"zero-admin/pkg/response/xerr"
	"zero-admin/rpc/sys/db/mysql/model"
	permLogic "zero-admin/rpc/sys/internal/logic/permissionservice"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取当前用户信息
func (l *GetCurrentUserLogic) GetCurrentUser(in *sysclient.GetCurrentUserRequest) (*sysclient.UserInfo, error) {
	userWithRoles, err := l.svcCtx.DB.GetUserWithRole(l.ctx, in.UserId)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, xerr.NewErrCode(xerr.ErrorUserNotExist)
	case err != nil:
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		return nil, xerr.NewErrCode(xerr.ErrorDb)
	}

	roleCodes := make([]string, 0, len(userWithRoles.Roles))
	for _, role := range userWithRoles.Roles {
		roleCodes = append(roleCodes, role.RoleCode)
	}
	// menus
	menus, _ := l.svcCtx.DB.GetMenusByRole(l.ctx, roleCodes)

	resp := &sysclient.UserInfo{
		Id:       userWithRoles.ID,
		Username: userWithRoles.Username,
		Status:   userWithRoles.Status,
		Roles:    buildRoles(userWithRoles.Roles),
		MenuTree: buildMenuTree(menus, 0),
		Email:    userWithRoles.Email,
		Mobile:   userWithRoles.Mobile,
		RealName: userWithRoles.RealName,
		Gender:   userWithRoles.Gender,
		Avatar:   userWithRoles.Avatar,
	}
	return resp, nil
}

func buildRoles(roles []model.SysRole) (res []*sysclient.Role) {
	res = make([]*sysclient.Role, 0, len(roles))
	for _, role := range roles {
		res = append(res, &sysclient.Role{
			Id:       role.ID,
			RoleName: role.RoleName,
			RoleCode: role.RoleCode,
		})
	}
	return
}

func buildMenuTree(menus []model.SysMenu, parentID int64) (menuTree []*sysclient.Menu) {
	menuTree = make([]*sysclient.Menu, 0)
	for _, menu := range menus {
		if menu.ParentID == parentID && menu.Status == 1 {
			m := permLogic.ConvertToRpcMenu(&menu)
			menuTree = append(menuTree, m)
			m.Children = buildMenuTree(menus, menu.ID)
		}
	}
	return
}
