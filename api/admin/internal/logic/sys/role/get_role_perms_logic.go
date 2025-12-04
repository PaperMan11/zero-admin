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

type GetRolePermsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermsLogic {
	return &GetRolePermsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolePermsLogic) GetRolePerms(req *types.IdValue) (resp *types.RoleInfo, err error) {
	res, err := l.svcCtx.RoleService.GetRolePerms(l.ctx, &roleservice.Int64Value{Value: req.Id})
	if err != nil {
		logc.Errorf(l.ctx, "获取角色权限失败: %v", err)
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
