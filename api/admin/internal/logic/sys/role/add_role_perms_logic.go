// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
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

// 废弃
func (l *AddRolePermsLogic) AddRolePerms(req *types.AddRolePermsRequest) (resp *types.RoleInfo, err error) {
	roleScopes := make([]*roleservice.RoleScope, 0, len(req.RoleScopes))
	for _, roleScope := range req.RoleScopes {
		roleScopes = append(roleScopes, &roleservice.RoleScope{
			Perms:     roleScope.Perms,
			RoleCode:  req.RoleCode,
			ScopeCode: roleScope.ScopeCode,
		})
	}
	res, err := l.svcCtx.RoleService.AddRolePerms(l.ctx, &roleservice.AddRolePermsRequest{
		RoleCode:   req.RoleCode,
		RoleScopes: roleScopes,
	})
	if err != nil {
		logc.Errorf(l.ctx, "添加角色权限失败: %v", err)
		return nil, err
	}

	// 添加casbin权限
	rules := make([][]string, 0, len(res.Scopes))
	for _, v := range res.Scopes {
		rules = append(rules, utils.ConvertToCasbinRule(res.Role.RoleCode, v.Scope.ScopeCode, v.Perms))
	}
	// AddPoliciesEx将授权规则添加到当前策略。 如果规则已经存在，规则将不会被添加。 但与AddPolicies不同，其他不存在的规则会被添加，而不是直接返回false
	ok, err := l.svcCtx.CasbinEnforcer.AddNamedPoliciesEx("p", rules)
	if err != nil || !ok {
		logc.Errorf(l.ctx, "添加casbin权限失败: %v", err)
	}

	scopes := make([]types.RoleScopeInfo, 0, len(res.Scopes))
	for _, v := range res.Scopes {
		scopes = append(scopes, types.RoleScopeInfo{
			Scope: utils.ConvertToTypesScope(v.Scope),
			Perms: v.Perms,
		})
	}
	return &types.RoleInfo{
		Role:   utils.ConvertToTypesRole(res.Role),
		Scopes: scopes,
	}, nil
}
