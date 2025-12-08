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

type DeleteRolePermsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRolePermsLogic {
	return &DeleteRolePermsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRolePermsLogic) DeleteRolePerms(req *types.DeleteRolePermsRequest) (resp *types.RoleInfo, err error) {
	rolePerms, err := l.svcCtx.RoleService.GetRolePerms(l.ctx, &roleservice.GetRolePermsRequest{RoleCode: req.RoleCode})
	if err != nil {
		logc.Errorf(l.ctx, "获取角色权限失败：%v", err)
		return nil, err
	}

	uid := utils.GetOperateID(l.ctx)
	res, err := l.svcCtx.RoleService.DeleteRolePerms(l.ctx, &roleservice.DeleteRolePermsRequest{
		RoleCode:   req.RoleCode,
		OperatorId: uid,
		ScopeCodes: req.ScopeCodes,
	})

	scopes := make([]types.RoleScopeInfo, 0, len(res.Scopes))
	for _, v := range res.Scopes {
		scopes = append(scopes, types.RoleScopeInfo{
			Scope: utils.ConvertToTypesScope(v.Scope),
			Perms: v.Perms,
		})
	}

	// 删除casbin权限
	rules := make([][]string, 0, len(req.ScopeCodes))
	for _, scopeCode := range req.ScopeCodes {
		for _, v := range rolePerms.Scopes {
			if v.Scope.ScopeCode == scopeCode {
				rules = append(rules, utils.ConvertToCasbinRule(res.Role.RoleCode, scopeCode, v.Perms))
			}
		}
	}
	ok, err := l.svcCtx.CasbinEnforcer.RemoveNamedPolicies("p", rules)
	if err != nil || !ok {
		logc.Errorf(l.ctx, "删除casbin权限失败: %v", err)
	}

	return &types.RoleInfo{
		Role:   utils.ConvertToTypesRole(res.Role),
		Scopes: scopes,
	}, nil
}
