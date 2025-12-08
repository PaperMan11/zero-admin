// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/sysclient"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleRoleStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleRoleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleRoleStatusLogic {
	return &ToggleRoleStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleRoleStatusLogic) ToggleRoleStatus(req *types.ToggleRoleStatusRequest) (resp *types.Role, err error) {
	uid := utils.GetOperateID(l.ctx)
	res, err := l.svcCtx.RoleService.ToggleRoleStatus(l.ctx, &sysclient.ToggleRoleStatusRequest{
		RoleCode:   req.RoleCode,
		Status:     req.Status,
		OperatorId: uid,
	})
	if err != nil {
		logc.Errorf(l.ctx, "切换角色状态失败: %v", err)
		return nil, err
	}

	// casbin权限
	var ok bool
	if req.Status == 1 {
		rolePerms, err := l.svcCtx.RoleService.GetRolePerms(l.ctx, &sysclient.GetRolePermsRequest{RoleCode: req.RoleCode})
		if err != nil {
			logc.Errorf(l.ctx, "获取角色权限失败: %v", err)
			return nil, err
		}
		// 添加casbin权限
		rules := make([][]string, 0, len(rolePerms.Scopes))
		for _, v := range rolePerms.Scopes {
			rules = append(rules, utils.ConvertToCasbinRule(rolePerms.Role.RoleCode, v.Scope.ScopeCode, v.Perms))
		}
		// AddPoliciesEx将授权规则添加到当前策略。 如果规则已经存在，规则将不会被添加。 但与AddPolicies不同，其他不存在的规则会被添加，而不是直接返回false
		ok, err = l.svcCtx.CasbinEnforcer.AddNamedPoliciesEx("p", rules)
		if err != nil || !ok {
			logc.Errorf(l.ctx, "添加casbin权限失败: %v", err)
		}
	} else {
		ok, err = l.svcCtx.CasbinEnforcer.RemoveFilteredNamedPolicy("p", 0, req.RoleCode)
		if err != nil || !ok {
			logc.Errorf(l.ctx, "删除角色权限失败: %v", err)
		}
	}

	return &types.Role{
		RoleId:      res.RoleId,
		RoleName:    res.RoleName,
		RoleCode:    res.RoleCode,
		Description: res.Description,
		Status:      res.Status,
	}, nil
}
