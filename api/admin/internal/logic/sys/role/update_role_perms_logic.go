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

func (l *UpdateRolePermsLogic) UpdateRolePerms(req *types.UpdateRolePermsRequest) (resp *types.RoleInfo, err error) {
	uid := logic.GetOperateID(l.ctx)
	roleScopes := make([]*roleservice.RoleScope, 0, len(req.RoleScopes))
	for _, v := range req.RoleScopes {
		roleScopes = append(roleScopes, &roleservice.RoleScope{
			Id:        v.Id,
			RoleCode:  v.RoleCode,
			ScopeCode: v.ScopeCode,
			Perms:     v.Perms,
		})
	}
	res, err := l.svcCtx.RoleService.UpdateRolePerms(l.ctx, &roleservice.UpdateRolePermsRequest{
		RoleId:     req.RoleId,
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
			Scope: logic.ConvertToTypesScope(v.Scope),
			Perms: v.Perms,
		})
	}
	return &types.RoleInfo{
		Role:   logic.ConvertToTypesRole(res.Role),
		Scopes: scopes,
	}, nil
}
