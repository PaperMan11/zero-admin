// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/rpc/sys/client/roleservice"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRolePermsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRolePermsLogic {
	return &AddRolePermsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRolePermsLogic) AddRolePerms(req *types.AddRolePermsRequest) (resp *types.RoleInfo, err error) {
	roleScopes := make([]*roleservice.RoleScope, 0, len(req.RoleScopes))
	for _, roleScope := range req.RoleScopes {
		roleScopes = append(roleScopes, &roleservice.RoleScope{
			Id:        roleScope.Id,
			Perms:     roleScope.Perms,
			RoleCode:  roleScope.RoleCode,
			ScopeCode: roleScope.ScopeCode,
		})
	}
	res, err := l.svcCtx.RoleService.AddRolePerms(l.ctx, &roleservice.AddRolePermsRequest{
		RoleId:     req.RoleId,
		RoleCode:   req.RoleCode,
		RoleScopes: roleScopes,
	})
	if err != nil {
		logc.Errorf(l.ctx, "添加角色权限失败: %v", err)
		return nil, err
	}

	scopes := make([]types.RoleScopeInfo, 0, len(res.Scopes))
	for _, v := range res.Scopes {
		scopes = append(scopes, types.RoleScopeInfo{
			Scope: types.Scope{
				Id:          v.Scope.Id,
				ScopeName:   v.Scope.ScopeName,
				ScopeCode:   v.Scope.ScopeCode,
				Description: v.Scope.Description,
				Sort:        v.Scope.Sort,
			},
			Perms: v.Perms,
		})
	}
	return &types.RoleInfo{
		Role: types.Role{
			RoleId:      res.Role.RoleId,
			RoleName:    res.Role.RoleName,
			RoleCode:    res.Role.RoleCode,
			Description: res.Role.Description,
			Status:      res.Role.Status,
		},
		Scopes: scopes,
	}, nil
}
