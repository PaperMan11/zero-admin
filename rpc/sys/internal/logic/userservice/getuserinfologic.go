package userservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"zero-admin/rpc/sys/db/common"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/logic"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取当前用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *sysclient.GetUserInfoRequest) (*sysclient.UserInfo, error) {
	user, err := l.svcCtx.DB.GetUserByID(l.ctx, in.UserId)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errors.New("用户不存在")
	case err != nil:
		logc.Errorf(l.ctx, "查询用户信息, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "查询用户信息异常")
	}

	// 用户角色信息
	var isSuperuser bool
	userRoles, _ := l.svcCtx.DB.GetRolesByUserID(l.ctx, user.ID)
	roleCodes := make([]string, 0, len(userRoles))
	for _, role := range userRoles {
		if common.IsSuperUser(role.RoleCode) {
			isSuperuser = true
		}
		roleCodes = append(roleCodes, role.RoleCode)
	}

	// 映射用户权限
	userPermMap := make(map[int64][]string)
	var menus []*model.SysMenu
	if isSuperuser {
		menus, _ = l.svcCtx.DB.GetAllMenus(l.ctx)
		for _, menu := range menus {
			userPermMap[menu.ScopeID] = common.PermissionMap[common.PERM_ALL]
		}
	} else {
		// menus
		menus, _ = l.svcCtx.DB.GetMenusByRoles(l.ctx, roleCodes)
		userPerms, _ := l.svcCtx.DB.GetRolesScopesPerm(l.ctx, roleCodes)
		for _, userPerm := range userPerms {
			userPermMap[userPerm.ID] = common.PermissionMap[userPerm.Perm]
		}
	}

	menuTree := logic.BuildMenuTree(menus, 0)
	// 添加菜单权限
	for _, menu := range menuTree {
		menu.Perms = userPermMap[menu.ScopeId]
	}

	resp := &sysclient.UserInfo{
		Id:       user.ID,
		Username: user.Username,
		Status:   user.Status,
		Roles:    logic.ConvertToRpcRoles(userRoles),
		MenuTree: menuTree,
		Email:    user.Email,
		Mobile:   user.Mobile,
		RealName: user.RealName,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
	}
	return resp, nil
}
