// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/logic"
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
			Scope: logic.ConvertToTypesScope(v.Scope),
			Perms: v.Perms,
		})
	}
	return &types.RoleInfo{
		Role:   logic.ConvertToTypesRole(res.Role),
		Scopes: scopes,
	}, nil
}
