// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"zero-admin/api/admin/internal/logic"
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
	uid := logic.GetOperateID(l.ctx)
	res, err := l.svcCtx.RoleService.DeleteRolePerms(l.ctx, &roleservice.DeleteRolePermsRequest{
		RoleId:     req.RoleId,
		RoleCode:   req.RoleCode,
		OperatorId: uid,
		ScopeCodes: req.ScopeCodes,
	})

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
