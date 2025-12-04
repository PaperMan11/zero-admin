package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-admin/rpc/sys/db/common"
	"zero-admin/rpc/sys/db/mysql/model"
	"zero-admin/rpc/sys/internal/svc"
	"zero-admin/rpc/sys/sysclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRolePermsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRolePermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRolePermsLogic {
	return &AddRolePermsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加角色权限
func (l *AddRolePermsLogic) AddRolePerms(in *sysclient.AddRolePermsRequest) (*sysclient.RoleInfo, error) {
	exists, err := l.svcCtx.DB.ExistsRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		logc.Errorf(l.ctx, "查询role_code失败, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "添加角色权限失败")
	}
	if !exists {
		logc.Errorf(l.ctx, "角色不存在, 参数：%+v", in)
		return nil, errors.New("角色不存在")
	}

	roleScopes := make([]*model.SysRoleScope, 0, len(in.GetRoleScopes()))
	for _, roleScope := range in.GetRoleScopes() {
		roleScopes = append(roleScopes, &model.SysRoleScope{
			RoleCode:  roleScope.RoleCode,
			ScopeCode: roleScope.ScopeCode,
			Perm:      common.ParsePermission(roleScope.Perms),
		})
	}

	err = l.svcCtx.DB.AddRoleScopes(l.ctx, roleScopes)
	if err != nil {
		logc.Errorf(l.ctx, "添加角色安全范围权限失败, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "添加角色安全范围权限失败")
	}

	perms, err := NewGetRolePermsLogic(l.ctx, l.svcCtx).GetRolePerms(&sysclient.Int64Value{Value: in.RoleId})
	if err != nil {
		logc.Errorf(l.ctx, "获取角色权限失败, 角色ID：%d, 异常: %s", in.RoleId, err.Error())
		return nil, status.Error(codes.Internal, "获取角色权限失败")
	}

	return perms, nil
}
