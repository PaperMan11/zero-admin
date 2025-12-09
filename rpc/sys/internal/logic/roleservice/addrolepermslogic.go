package roleservicelogic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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
	if common.IsSuperUser(in.RoleCode) {
		logc.Errorf(l.ctx, "超级管理员角色不允许修改, 角色：%s", in.RoleCode)
		return nil, status.Error(codes.PermissionDenied, common.ErrSuperUserDoNotEdit.Error())
	}
	//exists, err := l.svcCtx.DB.ExistsRoleByCode(l.ctx, in.RoleCode)
	role, err := l.svcCtx.DB.GetRoleByCode(l.ctx, in.RoleCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logc.Errorf(l.ctx, "角色不存在, 参数：%+v", in)
			return nil, errors.New("角色不存在")
		}
		logc.Errorf(l.ctx, "查询role_code失败, 参数：%+v, 异常: %s", in, err.Error())
		return nil, status.Error(codes.Internal, "添加角色权限失败")
	}

	scopesCodes := make([]string, 0, len(in.GetRoleScopes()))
	for _, roleScope := range in.GetRoleScopes() {
		scopesCodes = append(scopesCodes, roleScope.ScopeCode)
	}
	scopes, _ := l.svcCtx.DB.GetScopesByCodes(l.ctx, scopesCodes)
	if len(scopes) != len(scopesCodes) {
		logc.Errorf(l.ctx, "角色安全范围不存在, 参数：%+v", in)
		return nil, status.Error(codes.InvalidArgument, "角色安全范围不存在")
	}

	for _, roleScope := range in.GetRoleScopes() {
		err = l.svcCtx.DB.UpsertRoleScopes(l.ctx, model.SysRoleScope{
			RoleCode:  roleScope.RoleCode,
			ScopeCode: roleScope.ScopeCode,
			Perm:      common.ParsePermission(roleScope.Perms),
		})
		if err != nil {
			logc.Errorf(l.ctx, "添加角色安全范围权限失败, 参数：%+v, 异常: %s", in, err.Error())
			return nil, status.Error(codes.Internal, "添加角色安全范围权限失败")
		}
	}

	perms, err := NewGetRolePermsLogic(l.ctx, l.svcCtx).GetRolePerms(&sysclient.GetRolePermsRequest{RoleCode: role.RoleCode})
	if err != nil {
		logc.Errorf(l.ctx, "获取角色权限失败, 角色：%s, 异常: %s", role.RoleCode, err.Error())
		return nil, status.Error(codes.Internal, "获取角色权限失败")
	}

	return perms, nil
}
