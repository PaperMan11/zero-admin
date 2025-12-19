// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package scope

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"zero-admin/api/admin/internal/utils"
	"zero-admin/rpc/sys/client/roleservice"
	"zero-admin/rpc/sys/sysclient"

	"zero-admin/api/admin/internal/svc"
	"zero-admin/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleScopeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleScopeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleScopeStatusLogic {
	return &ToggleScopeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleScopeStatusLogic) ToggleScopeStatus(req *types.ToggleScopeStatusRequest) (resp *types.Scope, err error) {
	uid := utils.GetOperateID(l.ctx)
	res, err := l.svcCtx.ScopeService.ToggleScopeStatus(l.ctx, &sysclient.ToggleScopeStatusRequest{
		ScopeCode:  req.ScopeCode,
		Status:     req.Status,
		OperatorId: uid,
	})
	if err != nil {
		logc.Errorf(l.ctx, "切换安全范围状态失败: %v", err)
		return nil, err
	}

	// casbin权限
	if req.Status == 1 {
		rolesResp, _ := l.svcCtx.RoleService.GetRoleListByScopeCode(l.ctx, &roleservice.GetRolesByScopeCodeRequest{ScopeCode: req.ScopeCode})
		rules := make([][]string, 0, len(rolesResp.RolePerms))
		for _, v := range rolesResp.RolePerms {
			rules = append(rules, utils.ConvertToCasbinRule(v.RoleCode, v.ScopeCode, v.Perms))
		}
		// AddPoliciesEx将授权规则添加到当前策略。 如果规则已经存在，规则将不会被添加。 但与AddPolicies不同，其他不存在的规则会被添加，而不是直接返回false
		ok, err := l.svcCtx.CasbinEnforcer.AddNamedPoliciesEx("p", rules)
		if err != nil || !ok {
			logc.Errorf(l.ctx, "添加角色权限失败: %v, role_scope: %v", err, rolesResp.RolePerms)
		}
	} else {
		ok, err := l.svcCtx.CasbinEnforcer.RemoveFilteredNamedPolicy("p", 1, utils.ConvertScopeCodeToUrl(req.ScopeCode))
		if err != nil || !ok {
			logc.Errorf(l.ctx, "删除casbin权限失败: %v, scope: %s", err, req.ScopeCode)
		}
	}

	return &types.Scope{
		Id:          res.Id,
		ScopeName:   res.ScopeName,
		ScopeCode:   res.ScopeCode,
		Description: res.Description,
		Sort:        res.Sort,
	}, nil

}
