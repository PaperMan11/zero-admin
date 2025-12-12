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

type UpdateRolePermsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRolePermsLogic {
	return &UpdateRolePermsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 全量更新
func (l *UpdateRolePermsLogic) UpdateRolePerms(req *types.UpdateRolePermsRequest) (resp *types.RoleInfo, err error) {
	uid := utils.GetOperateID(l.ctx)
	roleScopes := make([]*roleservice.RoleScope, 0, len(req.RoleScopes))
	for _, v := range req.RoleScopes {
		roleScopes = append(roleScopes, &roleservice.RoleScope{
			RoleCode:  v.RoleCode,
			ScopeCode: v.ScopeCode,
			Perms:     v.Perms,
		})
	}
	res, err := l.svcCtx.RoleService.UpdateRolePerms(l.ctx, &roleservice.UpdateRolePermsRequest{
		RoleCode:   req.RoleCode,
		OperatorId: uid,
		RoleScopes: roleScopes,
	})
	if err != nil {
		logc.Errorf(l.ctx, "更新角色权限失败: %v", err)
		return nil, err
	}

	scopes := make([]types.RoleScopeInfo, 0, len(res.Scopes))
	for _, v := range res.Scopes {
		scopes = append(scopes, types.RoleScopeInfo{
			Scope: utils.ConvertToTypesScope(v.Scope),
			Perms: v.Perms,
		})
	}

	// 更新casbin
	l.svcCtx.CasbinEnforcer.RemoveFilteredNamedPolicy("p", 0, req.RoleCode)
	rules := make([][]string, 0, len(res.Scopes))
	for _, v := range res.Scopes {
		rules = append(rules, utils.ConvertToCasbinRule(res.Role.RoleCode, v.Scope.ScopeCode, v.Perms))
	}
	if len(rules) > 0 {
		ok, err := l.svcCtx.CasbinEnforcer.AddNamedPoliciesEx("p", rules)
		if err != nil || !ok {
			logc.Errorf(l.ctx, "更新casbin权限失败: %v", err)
		}
	}

	return &types.RoleInfo{
		Role:   utils.ConvertToTypesRole(res.Role),
		Scopes: scopes,
	}, nil
}
